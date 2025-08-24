import json
import logging
from elastic_search import add_to_elastic_search

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def callback(ch, method, properties, body):
    try:
        # Decode JSON message
        log_entry = json.loads(body.decode("utf-8"))
        logger.info(f"Processing log message: {log_entry.get('id', 'unknown ID')}")

        add_to_elastic_search(log_entry)
        # Acknowledge message after successful processing
        ch.basic_ack(delivery_tag=method.delivery_tag)
        logger.info(f"Successfully processed and acknowledged message: {log_entry.get('id', 'unknown ID')}")

    except json.JSONDecodeError as e:
        logger.error(f"JSON decode error: {e}. Message body: {body}")
        # Reject malformed messages without requeueing
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    except Exception as e:
        logger.error(f"Error processing message: {e}")
        # Reject message without requeueing to prevent infinite loops
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

#  [x] Received: {'id': 'd98e7a5e-530e-42aa-b213-fdf02a9e984b', 'time_stamp': '2025-08-22T16:56:18.899882Z', 'service': 'Authentication', 'level': 'Info', 'message': 'User Login Sucessfully', 'meta_data': {'source_ip': '127.0.0.1', 'region': 'asia-south-01'}}