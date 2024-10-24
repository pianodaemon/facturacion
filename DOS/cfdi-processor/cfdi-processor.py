import time
import redis
import os


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
    return redis.StrictRedis.from_url(connection_str)


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

    def __call__(self, interval):
        """Process messages in a loop"""
        while True:
            message = self.input_queue.pop()
            if message:
                self.handle_message(message)
            time.sleep(interval)

    @staticmethod
    def handle_message(message):
        """Handle processing of a message"""
        print("Processing message:", message)
        
        # Example: use the message to extract relevant data for the payload
        message_data = json.loads(message)
        
        url = "{ HOST }/v4/cfdi40/create"
        payload = json.dumps({
            "Receptor": {
                "UID": message_data.get("UID", "6169fc02637e1")
            },
            "TipoDocumento": "factura",
            "Conceptos": [
                {
                    "ClaveProdServ": "81112101",
                    "Cantidad": 1,
                    "ClaveUnidad": "E48",
                    "Unidad": "Unidad de servicio",
                    "ValorUnitario": 229.9,
                    "Descripcion": "Desarrollo a la medida",
                    "Impuestos": {
                        "Traslados": [
                            {
                                "Base": 229.9,
                                "Impuesto": "002",
                                "TipoFactor": "Tasa",
                                "TasaOCuota": "0.16",
                                "Importe": 36.784
                            }
                        ],
                        "Locales": [
                            {
                                "Base": 229.9,
                                "Impuesto": "ISH",
                                "TipoFactor": "Tasa",
                                "TasaOCuota": "0.03",
                                "Importe": 6.897
                            }
                        ]
                    }
                }
            ],
            "UsoCFDI": "P01",
            "Serie": 17317,
            "FormaPago": "03",
            "MetodoPago": "PUE",
            "Moneda": "MXN",
            "EnviarCorreo": False
        })
        headers = {
            'Content-Type': 'application/json',
            'F-PLUGIN': '9d4095c8f7ed5785cb14c0e3b033eeb8252416ed',
            'F-Api-Key': 'Tu API key',
            'F-Secret-Key': 'Tu Secret key'
        }

        try:
            response = requests.post(url, headers=headers, data=payload)
            print("API response:", response.text)
        except requests.exceptions.RequestException as e:
            print(f"API request failed: {e}")

if __name__ == "__main__":
    processor = Processor("billingInput", "billingNotifications")
    processor(10)
