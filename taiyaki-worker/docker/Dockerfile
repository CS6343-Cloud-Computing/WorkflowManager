# Base DockerFile
FROM cs6343g7/cleanarticle:finalterm
ENV user_cmd "python3 ./cleanArticle.py"
ENV consume_topic a79417ab-bd83-4874-b643-c9fbdebc4888
ENV kafka_server "192.168.1.82:9092"

RUN apt-get update
RUN apt-get install -y python3 python3-pip

RUN pip3 install kafka-python

COPY ./streamExecutor.py ./

CMD ["python3", "./streamExecutor.py"]
