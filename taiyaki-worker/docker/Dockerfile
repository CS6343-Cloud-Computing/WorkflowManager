# Base DockerFile
FROM cs6343g7/classifier:finalterm
ENV user_cmd "python3 ./predict.py"
ENV consume_topic 0e309204-05dc-4562-86ad-46e9c807f12c
ENV PYTHONUNBUFFERED=2
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3","-u", "./streamExecutor.py"]
