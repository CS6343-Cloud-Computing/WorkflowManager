from kafka import KafkaConsumer
from kafka import KafkaProducer
from kafka.structs import TopicPartition
import os
import subprocess


class streamExecutor:
    def __init__(self):
        self.kafkaServer =  os.getenv('kafka_server')
        self.consumeTopic = os.getenv('consume_topic')
        self.userCMD = os.getenv('user_cmd')
        self.kafka_init()

    def kafka_init(self):

        self.producer = KafkaProducer(
            bootstrap_servers=self.kafkaServer, api_version=(0, 11, 5))
        
        while True:

            try:
                self.consumer = KafkaConsumer(
                    bootstrap_servers=self.kafkaServer, api_version=(0, 11, 5), auto_offset_reset="earliest", group_id=None
                )
                
                self.consumer.assign([TopicPartition(self.consumeTopic, 0)])
                break
            except:
                print("Consumer not ready!!")

    def Exec(self):
        for msg in self.consumer:
            data = msg.value
            print(data)
            cmds = self.userCMD.split()
            cmds.append(data)
            userProcess = subprocess.Popen(cmds, stdout=subprocess.PIPE, text=True, universal_newlines=True)
                                           
            output = str(userProcess.communicate()[0])
            print(msg)
            print(output, end='')
            flow = str(msg.headers[0][1])
            flow = flow.split("<===>")
            pointer = int(msg.headers[1][1])

            self.producer.send(flow[pointer],  key = msg.key,  value = msg.value,partition=0, headers=[('flow', "<===>".join(flow).encode('utf-8')), ('pointer', str(pointer+1).encode())])

            


if __name__ == "__main__":
    se = streamExecutor()
    se.Exec()
