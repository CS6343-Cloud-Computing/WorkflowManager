# Base DockerFile
FROM ubuntu
ENV user_cmd "sudo test"
ENV consume_topic c329a42c-5c3a-4d39-97d9-23e4bd551fd3
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
