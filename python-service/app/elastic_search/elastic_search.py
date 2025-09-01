import os
from elasticsearch import Elasticsearch


# Use localhost for local development, elasticsearch for Docker
ES_URL = os.getenv("ELASTICSEARCH_URL", "http://localhost:9200")

# Global variable to store ES client
es = None
index_name = "logs"

mapping = {
    "mappings": {
        "properties": {
            "id": {"type": "keyword"},
            "time_stamp": {"type": "date", "format": "strict_date_optional_time||epoch_millis"},
            "service": {"type": "keyword"},
            "level": {"type": "keyword"},
            "message": {"type": "text"},
            "meta_data": {
                "properties": {
                    "source_ip": {"type": "ip"},
                    "region": {"type": "keyword"}
                }
            }
        }
    }
}

def get_elasticsearch_client():
    """Get or create Elasticsearch client with connection handling."""
    global es
    if es is None:
        try:
            es = Elasticsearch(ES_URL)
            # Test the connection
            es.info()
            print(f"Successfully connected to Elasticsearch at {ES_URL}")
            
            # Create index if it doesn't exist
            if not es.indices.exists(index=index_name):
                es.indices.create(index=index_name, body=mapping)
                print(f"Created Elasticsearch index: {index_name}")
        except Exception as e:
            print(f"Failed to connect to Elasticsearch at {ES_URL}: {e}")
            print("Make sure Elasticsearch is running on localhost:9200 or set ELASTICSEARCH_URL environment variable")
            es = None
    return es


def add_to_elastic_search(process_data: dict) -> bool:
    """
    Pushes a document to Elasticsearch.
    
    Args:
        process_data (dict): The document to insert.
    
    Returns:
        bool: True if successful, False otherwise.
    """
    try:
        client = get_elasticsearch_client()
        if client is None:
            print("Elasticsearch client not available")
            return False
            
        client.index(index=index_name, document=process_data)
        print(f"Document pushed to Elasticsearch index: {index_name}")
        return True
    except Exception as e:
        print("Exception occurred while pushing to Elasticsearch:", e)
        return False
