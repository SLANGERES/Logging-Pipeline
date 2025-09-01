import json

from elastic_search.elastic_search import add_to_elastic_search



def callback(ch, method, properties, body):
    try:
        # Decode JSON message
        log_entry = json.loads(body.decode("utf-8"))
        print(f"Processing log message: {log_entry.get('id', 'unknown ID')}")

        add_to_elastic_search(log_entry)
        # Acknowledge message after successful processing
        ch.basic_ack(delivery_tag=method.delivery_tag)
        print(f"Successfully processed and acknowledged message: {log_entry.get('id', 'unknown ID')}")

    except json.JSONDecodeError as e:
        print(f"JSON decode error: {e}. Message body: {body}")
        # Reject malformed messages without requeueing
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    except Exception as e:
        print(f"Error processing message: {e}")
        # Reject message without requeueing to prevent infinite loops
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

