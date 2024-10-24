import time
import redis
import json

# Connect to Redis
r = redis.StrictRedis(host='localhost', port=6379, db=0)

def do_processing():
    while True:

        # Sleep for a bit before checking again
        time.sleep(10)  # Check every 10 seconds


if __name__ == "__main__":
    do_processing()
