# Base DockerFile
FROM cs6343g7/sentimentanalyzer:midterm
ENV user_cmd "sudo test"
ENV consume_topic e82b6d7b-d91c-4ef2-8c59-106bea790521
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
