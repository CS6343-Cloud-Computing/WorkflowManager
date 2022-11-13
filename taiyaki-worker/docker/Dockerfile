# Base DockerFile
FROM ubuntu
ENV user_cmd "sudo test"
ENV consume_topic 17c2e4b9-bd61-42c0-a69d-9281f36f4d2e
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
