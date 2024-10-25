import time
import redis
import traceback
import os
import sys
import requests
import json


def redis_connected(func):
    """Decorator to manage Redis connection for a function"""

    def wrapper(*args, **kwargs):
        client = connect_to_redis()
        try:
            return func(client, *args, **kwargs)
        except Exception as e:
            raise e
        finally:
            client.close()

    return wrapper


def connect_to_redis():
    """Opens a connection to Redis"""
    redis_host = os.getenv("REDIS_HOST")
    redis_port = os.getenv("REDIS_PORT")

    if not redis_host or not redis_port:
        raise ValueError("Environment variables REDIS_HOST or REDIS_PORT not set")

    connection_str = f"redis://{redis_host}:{redis_port}/0"
    return redis.Redis.from_url(connection_str)  # Updated to use Redis instead of StrictRedis


class QueueNotFoundException(Exception):
    """Custom exception for missing Redis queues."""
    pass


class RedisQueue:
    """Redis queue handler"""

    def __init__(self, queue_name):
        self.queue_name = queue_name

    @redis_connected
    def is_present(self, client):
        return client.exists(self.queue_name)

    @redis_connected
    def ping(self, client):
        """Check the connection to Redis."""
        client.ping()

    @redis_connected
    def push(self, client, message):
        """Push a message to the tail of the queue"""
        client.lpush(self.queue_name, message)

    @redis_connected
    def pop(self, client):
        """Pop a message from the head of the queue (blocking)"""
        m = client.brpop(self.queue_name)
        return m[1] if m else None


class Processor:
    """Processor for handling Redis queues"""

    def __init__(self, input_queue_name, notification_queue_name):
        self.input_queue = RedisQueue(input_queue_name)
        self.notification_queue = RedisQueue(notification_queue_name)

    def __call__(self, shaper_handler, dispatcher_handler, interval):
        """Process messages in a loop"""
        if not self.input_queue.is_present():
            raise QueueNotFoundException(f"Input queue '{self.input_queue.queue_name}' does not exist yet.")

        while True:
            message = self.input_queue.pop()
            if message:
                self.handle_message(shaper_handler, dispatcher_handler, message)
            time.sleep(interval)

    @staticmethod
    def handle_message(shaper_handler, dispatcher_handler, message):  # Fixed the syntax here
        """Handle processing of a message"""
        print("Processing message:", message)

        # Parse the message, which is in JSON format (representing a Receipt object)
        receipt = json.loads(message)
        dispatcher_handler(shaper_handler(receipt))


def factura_com_shaper(receipt):
    """Construct the payload for the third-party API"""
    return {
        "Receptor": {
            "UID": receipt.get("receptor_data_ref")
        },
        "TipoDocumento": "factura",
        "Conceptos": [
            {
                "ClaveProdServ": item.get("fiscal_product_id"),
                "Cantidad": item.get("product_quantity"),
                "ClaveUnidad": item.get("fiscal_product_unit"),
                "Unidad": item.get("product_unidad"),
                "ValorUnitario": item.get("product_unit_price"),
                "Descripcion": item.get("product_desc"),
                "Impuestos": {
                    "Traslados": [
                        {
                            "Base": transfer.get("base"),
                            "Impuesto": transfer.get("fiscal_type"),
                            "TipoFactor": transfer.get("fiscal_factor"),
                            "TasaOCuota": str(transfer.get("rate")),
                            "Importe": transfer.get("amount")
                        }
                        for transfer in item.get("product_transfers", [])
                    ]
                }
            }
            for item in receipt.get("items", [])
        ],
        "UsoCFDI": receipt.get("purpose"),
        "Serie": 17317,  # Assuming static, customize if needed
        "FormaPago": receipt.get("payment_way"),
        "MetodoPago": receipt.get("payment_method"),
        "Moneda": receipt.get("document_currency"),
        "TipoCambio": str(receipt.get("exchange_rate")),
        "EnviarCorreo": False  # Assuming no email, customize if needed
    }


def factura_com_dispatcher(payload):
    """Send the payload to the third-party API"""

    headers = {
        'Content-Type': 'application/json',
        'F-PLUGIN': '9d4095c8f7ed5785cb14c0e3b033eeb8252416ed',
        'F-Api-Key': 'Your API key',
        'F-Secret-Key': 'Your Secret key'
    }

    try:
        # Replace HOST with the actual API endpoint
        response = requests.post("{HOST}/v4/cfdi40/create", headers=headers, data=json.dumps(payload))
        print("API response:", response.text)
    except requests.exceptions.RequestException as e:
        print(f"API request failed: {e}")


if __name__ == "__main__":
    debug = os.getenv("DEBUG", "False").lower() == "true"  # Get debug flag from environment
    input_queue_name = os.getenv("INPUT_QUEUE_NAME", "billingInput")
    notification_queue_name = os.getenv("NOTIFICATION_QUEUE_NAME", "billingNotifications")
    try:
        processor = Processor("billingInput", "billingNotifications")
        processor(factura_com_shaper, factura_com_dispatcher, 10)
    except KeyboardInterrupt:
        print('Exiting')
    except Exception as e:
        if debug:
            print(f'Whoops! problem in processor: {e}', file=sys.stderr)
            traceback.print_exc(file=sys.stderr)
        sys.exit(1)
