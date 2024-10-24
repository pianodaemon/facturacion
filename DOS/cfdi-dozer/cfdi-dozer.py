import time
import redis
import json

# Connect to Redis
r = redis.StrictRedis(host='localhost', port=6379, db=0)

def do_polling():
    while True:
        # Get the current timestamp
        current_time = int(time.time())

        # Fetch bills whose scheduled time is due
        scheduled_bills = r.zrangebyscore("billingBacklog", 0, current_time)

        if scheduled_bills:
            for bill in scheduled_bills:
                # Send the email
                forward(bill)

                # Remove the bill from Redis after sending
                r.zrem("billingBacklog", bill)

        # Sleep for a bit before checking again
        time.sleep(10)  # Check every 10 seconds

def forward(bill):
    pass

if __name__ == "__main__":
    do_polling()
