import os
import sys
from Rabbit_mq.mq_client import callback 
from Rabbit_mq.mq_channel import Mq_connection


def main():
    try:
        # Step 1: Connect to RabbitMQ
        # Use environment variable or default to localhost for local development
        rabbit_host = os.getenv("RABBITMQ_HOST", "localhost")
        channel = Mq_connection(rabbit_host)
        
        if channel is None:
            print("Failed to establish RabbitMQ connection. Exiting.")
            sys.exit(1)

        # Step 2: Declare queue (must match producer)
        queue_name = "logs"
        channel.queue_declare(queue=queue_name, durable=True)
        print(f"Declared queue: {queue_name}")

        # Step 3: Start consuming
        channel.basic_consume(
            queue=queue_name,
            on_message_callback=callback,
            auto_ack=False  # manual ack
        )

        print(" [*] Waiting for messages. To exit press CTRL+C")
        channel.start_consuming()
        
    except KeyboardInterrupt:
        print("\n [!] Interrupted by user")
        if 'channel' in locals():
            channel.stop_consuming()
            channel.connection.close()
        sys.exit(0)
    except Exception as e:
        print(f"Unexpected error in main: {e}")
        sys.exit(1)
        
if __name__ == "__main__":
    main()
