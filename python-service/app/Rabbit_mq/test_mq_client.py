import json
import pytest
from unittest.mock import Mock, patch
from Rabbit_mq import callback  # assuming your code is in consumer.py

# --------------- SUCCESS CASE ----------------
@patch("consumer.add_to_elastic_search")  # mock elasticsearch
def test_callback_success(mock_add_to_es):
    # Mock RabbitMQ channel + method
    mock_channel = Mock()
    mock_method = Mock()
    mock_method.delivery_tag = "test_tag"

    # A valid log message
    log_entry = {"id": "123", "message": "hello"}
    body = json.dumps(log_entry).encode("utf-8")

    callback(mock_channel, mock_method, None, body)

    # ✅ add_to_elastic_search should be called with parsed dict
    mock_add_to_es.assert_called_once_with(log_entry)

    # ✅ basic_ack should be called
    mock_channel.basic_ack.assert_called_once_with(delivery_tag="test_tag")

    # ❌ basic_nack should not be called
    mock_channel.basic_nack.assert_not_called()


# --------------- JSON DECODE ERROR CASE ----------------
def test_callback_json_decode_error():
    mock_channel = Mock()
    mock_method = Mock()
    mock_method.delivery_tag = "bad_json"

    # Invalid JSON body
    body = b"{invalid-json}"

    callback(mock_channel, mock_method, None, body)

    # ✅ Should call basic_nack with requeue=False
    mock_channel.basic_nack.assert_called_once_with(delivery_tag="bad_json", requeue=False)
    mock_channel.basic_ack.assert_not_called()


# --------------- ELASTICSEARCH FAILURE CASE ----------------
@patch("consumer.add_to_elastic_search", side_effect=Exception("Elasticsearch down"))
def test_callback_elasticsearch_failure(mock_add_to_es):
    mock_channel = Mock()
    mock_method = Mock()
    mock_method.delivery_tag = "es_fail"

    body = json.dumps({"id": "999"}).encode("utf-8")

    callback(mock_channel, mock_method, None, body)

    # ✅ add_to_elastic_search should be attempted
    mock_add_to_es.assert_called_once()

    # ✅ Should nack since ES failed
    mock_channel.basic_nack.assert_called_once_with(delivery_tag="es_fail", requeue=False)
    mock_channel.basic_ack.assert_not_called()
