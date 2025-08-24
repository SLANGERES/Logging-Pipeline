import pika
from Rabbit_mq import callback 

def main():
    # Step 1: Connect to RabbitMQ
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="localhost")
    )
    channel = connection.channel()

    # Step 2: Declare queue (must match producer)
    queue_name = "logs"
    channel.queue_declare(queue=queue_name, durable=True)

    # Step 3: Start consuming
    channel.basic_consume(
        queue=queue_name,
        on_message_callback=callback,
        auto_ack=False  # manual ack
    )

    print(" [*] Waiting for messages. To exit press CTRL+C")
    channel.start_consuming() 
if __name__ == "__main__":
    main()
