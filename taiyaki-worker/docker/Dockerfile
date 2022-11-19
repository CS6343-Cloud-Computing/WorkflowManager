# Base DockerFile
FROM cs6343g7/sentimentanalyzer:finalterm
ENV user_cmd "python3 ./sentimentAnalyser.py"
ENV consume_topic 938d2d16-3201-4fa0-8511-79074a7eb3b1
ENV PYTHONUNBUFFERED=2
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3","-u", "./streamExecutor.py"]
