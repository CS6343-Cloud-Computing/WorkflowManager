from kafka import KafkaConsumer
from kafka import KafkaProducer
from kafka.structs import TopicPartition
import os
import subprocess
from collections import defaultdict
import queue

class MultiQueue:
    def __init__(self):
        self.indegree = os.getenv('indegree')
        self.queues = [queue.Queue()]*self.indegree
    
    def IsReady(self):
        # iterate over all queues and check the top element
        # if present then only return true otherwise false

    def GetList(self):
        # returns the first element of the each queue as a list



class streamExecutor:
    def __init__(self):
        self.kafkaServer =  os.getenv('kafka_server')
        self.consumeTopic = os.getenv('consume_topic')
        self.userCMD = os.getenv('user_cmd')
        self.streamExecutor_init()
        self.kafka_init()
    
    def streamExecutor_init(self):
        self.messageQueue = defaultdict()

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
        print("Started stream executor")
        for msg in self.consumer:
            # Get the workflow ID of the message
            # Get the sender ID of the message
            # Insert the message inside the queue for that particular senderID queue in the workflow
            # Check the top of the queue for every queue in the workflow for which we just received the message
            # if the top of the queue is empty for some then wait
            # Otherwise call the user function and produce a next message

            data = msg.value.decode("utf-8")
            print(data)
            cmds = self.userCMD.split()
            cmds.append(data)
            userProcess = subprocess.Popen(cmds, stdout=subprocess.PIPE, text=True, universal_newlines=True)
                                           
            output = str(userProcess.communicate()[0])
            print(msg)
            print(output, end='')
            flow = msg.headers[0][1].decode("utf-8")
            flow = flow.split("<===>")
            pointer = int(msg.headers[1][1])
            workflow = msg.headers[2][1].decode("utf-8")
            if pointer >= len(flow):
                self.producer.send(workflow+"-output",  key = msg.key,  value = output.encode('utf-8'),partition=0, headers=[('flow', "<===>".join(flow).encode('utf-8')), ('pointer', str(pointer+1).encode('utf-8')),('workflow', workflow.encode('utf-8'))])
                continue

            self.producer.send(flow[pointer],  key = msg.key,  value = output.encode('utf-8'),partition=0, headers=[('flow', "<===>".join(flow).encode('utf-8')), ('pointer', str(pointer+1).encode('utf-8')),('workflow', workflow.encode('utf-8'))])

            
if __name__ == "__main__":
    se = streamExecutor()
    se.Exec()
