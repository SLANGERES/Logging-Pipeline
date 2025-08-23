import json
from elastic_search import add_to_elastic_search

def callback(ch, method, properties, body):
    try:
        # Decode JSON message
        log_entry = json.loads(body.decode("utf-8"))

        add_to_elastic_search(log_entry)
        # Acknowledge message after successful processing
        ch.basic_ack(delivery_tag=method.delivery_tag)

    except Exception as e:
        print(" [!] Error processing message:", e)
        # Reject message without requeueing
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

#  [x] Received: {'id': 'd98e7a5e-530e-42aa-b213-fdf02a9e984b', 'time_stamp': '2025-08-22T16:56:18.899882Z', 'service': 'Authentication', 'level': 'Info', 'message': 'User Login Sucessfully', 'meta_data': {'source_ip': '127.0.0.1', 'region': 'asia-south-01'}}