import time
import redis
import json


def redis_connected(func):
    """Handy decorator to fetch a database connection"""

    def wrapper(sql):
        c = _connect()
        try:
            return func(c, sql)
        except:
            raise
        finally:
            c.close()

    return wrapper

def _connect():
    """Opens a connection to database"""
    # order here matters
    env_vars = (
        "REDIS_HOST",
        "REDIS_PORT",
    )

    t = tuple(map(env_property, env_vars))
    try:
        conn_str = "redis://{0}:{1}/0".format(*t)
        return redis.StrictRedis.from_url(conn_str)
    except:
        raise

class QueueRedis(object):

    def __init__(self, queue_name):
        self._queue_name = queue_name

    @redis_connected
    def ping(self, client):
        """Check the connection otherwise raise exception."""
        client.ping()

    @redis_connected
    def push(self, client, msg):
        """Push to tail of the queue"""
        client.lpush(self._queue_name, msg)

    @redis_connected
    def pop(self, client):
        """Pop from head of the queue, It's blocking)"""
        _, msg = client.brpop(self._queue_name)
        return msg


def do_processing(input_queue, notifications_queue, span):
    """Start the processing cycle"""
    while True:
        msg = input_queue.pop()
        do_bill(msg)
        # Sleep for a bit before checking again
        time.sleep(span)  # Check every n seconds


def do_bill(msg):
    """Submit data to pac for billing issue"""
    print("Here we're gonna incept a factura")
    pass


if __name__ == "__main__":
    input_queue, notifications_queue = QueueRedis("billingInput"), QueueRedis("billingNotifications")
    do_processing(input_queue, notifications_queue, 10)
