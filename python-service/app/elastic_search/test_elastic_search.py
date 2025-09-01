import pytest
import requests
from datetime import datetime
from elastic_search import add_to_elastic_search



test_data={
            "id": "test-01",
            "time_stamp": datetime.time(),
            "service": "test-service",
            "level": "test",
            "message": "testing purpose",
            "meta_data": {
                "properties": {
                    "source_ip": "0.0.0.0",
                    "region": "test-region"
                }
            }
}   
def test_elastic_search():
    assert add_to_elastic_search(test_data)==True

def test_connection():
    url = "http://127.0.0.1:9200"   
    response = requests.get(url)
    assert response.status_code == 200