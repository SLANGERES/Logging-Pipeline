import os
import pika 
import time

def Mq_connection(host: str | None = None, max_retries: int = 5):
    # Use localhost for local development, rabbitmq for Docker
    host = host or os.getenv("RABBITMQ_HOST", "localhost")
    
    # Use port 5673 for local development (mapped from Docker), 5672 for Docker network
    port = 5673 if host == "localhost" else 5672
    
    for attempt in range(max_retries):
        try:
            print(f"Attempting to connect to RabbitMQ at {host}:{port} (attempt {attempt + 1}/{max_retries})")
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(host=host, port=port)
            )
            channel = connection.channel()
            print(f"Successfully connected to RabbitMQ at {host}")
            return channel
        except pika.exceptions.AMQPConnectionError as e:
            print(f"Failed to connect to RabbitMQ at {host}: {e}")
            if attempt < max_retries - 1:
                print(f"Retrying in 2 seconds...")
                time.sleep(2)
            else:
                print(f"Failed to connect to RabbitMQ after {max_retries} attempts")
                if host == "localhost":
                    print("Make sure RabbitMQ is running. You can start it with: docker-compose up rabbitmq")
                    print("Or install RabbitMQ locally and ensure it's running on port 5673")
                else:
                    print(f"Make sure RabbitMQ is running at {host}:{port}")
                return None
        except Exception as e:
            print(f"Unexpected error connecting to RabbitMQ: {e}")
            return None
        
