
from elasticsearch import Elasticsearch


es = Elasticsearch("http://localhost:9200") 


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

if not es.indices.exists(index=index_name):
    
    es.indices.create(index=index_name, body=mapping)

def add_to_elastic_search(process_data:dict):

   

    es.index(index=index_name,document=process_data)
    print("push to elastic serach ")