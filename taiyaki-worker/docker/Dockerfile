# Base DockerFile
FROM alpine
ENV user_cmd "sudo test"
ENV consume_topic e1651144-326a-45d5-9cc7-c972d4972082
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
