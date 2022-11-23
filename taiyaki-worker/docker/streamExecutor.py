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
        for queue in self.queues:
            if queue.empty():
                return False
        return True

    def GetList(self):
        result = []
        for queue in self.queues:
            result.append(queue.get())
        return result

    def PutElement(self, item, k):
        self.queues[k].put(item)



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
            workflow = msg.headers[2][1].decode("utf-8")
            sourceID = msg.headers[3][1].decode("utf-8")
            if self.messageQueue[workflow] is None:
                self.messageQueue[workflow] = MultiQueue()
            self.messageQueue[workflow].PutElement(msg, sourceID)
            if self.messageQueue[workflow].IsReady():
                res = self.messageQueue[workflow].GetList()
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

                if pointer >= len(flow):
                    self.producer.send(workflow+"-output",  key = msg.key,  value = output.encode('utf-8'),partition=0, headers=[('flow', "<===>".join(flow).encode('utf-8')), ('pointer', str(pointer+1).encode('utf-8')),('workflow', workflow.encode('utf-8'))])
                    continue

                self.producer.send(flow[pointer],  key = msg.key,  value = output.encode('utf-8'),partition=0, headers=[('flow', "<===>".join(flow).encode('utf-8')), ('pointer', str(pointer+1).encode('utf-8')),('workflow', workflow.encode('utf-8'))])
                # Otherwise call the user function and produce a next message

if __name__ == "__main__":
    se = streamExecutor()
    se.Exec()
