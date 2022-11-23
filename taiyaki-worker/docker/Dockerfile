# Base DockerFile
FROM cs6343g7/wf2c3
ENV user_cmd "python3 ./minMaxTemp.py"
ENV consume_topic 4c9e9802-0f3f-485e-be0b-39d348c04f5c
ENV PYTHONUNBUFFERED=2
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3","-u", "./streamExecutor.py"]
