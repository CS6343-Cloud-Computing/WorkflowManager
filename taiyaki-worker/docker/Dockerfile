# Base DockerFile
FROM alpine
ENV user_cmd "sudo test"
ENV consume_topic 3cc2809e-31cc-4a22-a591-6ae2da08ac9e
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
