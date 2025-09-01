#!/usr/bin/env python3
"""
Test script to verify connection fixes
"""

import sys
import os

# Add the app directory to the Python path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'app'))

def test_elasticsearch_connection():
    """Test Elasticsearch connection"""
    print("Testing Elasticsearch connection...")
    try:
        from elastic_search.elastic_search import get_elasticsearch_client
        client = get_elasticsearch_client()
        if client:
            print("✅ Elasticsearch connection successful!")
            return True
        else:
            print("❌ Elasticsearch connection failed!")
            return False
    except Exception as e:
        print(f"❌ Elasticsearch connection error: {e}")
        return False

def test_rabbitmq_connection():
    """Test RabbitMQ connection"""
    print("\nTesting RabbitMQ connection...")
    try:
        from Rabbit_mq.mq_channel import Mq_connection
        channel = Mq_connection()
        if channel:
            print("✅ RabbitMQ connection successful!")
            # Clean up
            channel.connection.close()
            return True
        else:
            print("❌ RabbitMQ connection failed (this is expected if RabbitMQ is not running)")
            return False
    except Exception as e:
        print(f"❌ RabbitMQ connection error: {e}")
        return False

def test_imports():
    """Test that all imports work correctly"""
    print("\nTesting imports...")
    try:
        from Rabbit_mq.mq_client import callback
        from Rabbit_mq.mq_channel import Mq_connection
        from elastic_search.elastic_search import add_to_elastic_search
        print("✅ All imports successful!")
        return True
    except Exception as e:
        print(f"❌ Import error: {e}")
        return False

def main():
    print("Connection Test Script")
    print("=" * 50)
    
    # Test imports first
    imports_ok = test_imports()
    
    # Test connections
    es_ok = test_elasticsearch_connection()
    rabbitmq_ok = test_rabbitmq_connection()
    
    print("\n" + "=" * 50)
    print("Summary:")
    print(f"Imports: {'✅' if imports_ok else '❌'}")
    print(f"Elasticsearch: {'✅' if es_ok else '❌'}")
    print(f"RabbitMQ: {'✅' if rabbitmq_ok else '❌ (may be expected if not running)'}")
    
    if imports_ok and es_ok:
        print("\n🎉 Core functionality is working! The Python service should now run properly.")
        print("\nTo start the service:")
        print("1. Make sure RabbitMQ is running on localhost:5672")
        print("2. Run: python app/main.py")
    else:
        print("\n⚠️  Some issues remain. Check the error messages above.")

if __name__ == "__main__":
    main()
