# Base DockerFile
FROM cs6343g7/sentimentanalyzer:finalterm
ENV user_cmd "python3 ./sentimentAnalyser.py"
ENV consume_topic 85dd4132-06ba-4cde-a0f3-5c2bc22bb7b9
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
