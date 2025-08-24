FROM python:3.12

WORKDIR /app

# Install dependencies
COPY requirements.txt ./
RUN pip install -r requirements.txt

# Copy source
COPY . .

# Run main.py inside /app/app/
CMD ["python3", "/app/app/main.py"]
