from kafka import KafkaConsumer
from kafka import KafkaProducer
from kafka.structs import TopicPartition
import os
import subprocess
from collections import defaultdict
import queue

class MultiQueue:
    def __init__(self):
        self.indegree = int(os.getenv('indegree'))
        self.queues = [queue.Queue() for i in range(self.indegree)]
    
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
        self.messageQueue = defaultdict(lambda: None)

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
            print("Rec msg: ", msg)
            input = eval(msg.headers[0][1].decode("utf-8"))
            output = eval(msg.headers[1][1].decode("utf-8"))
            workflow = msg.headers[2][1].decode("utf-8")
            containerName = os.getenv("consume_topic")
            if self.messageQueue[workflow] is None:
                self.messageQueue[workflow] = MultiQueue()
            prev = msg.key.decode("utf-8").split("::")[1]
            prevIndex = input[containerName].index(prev)
            print("inserting in the queue: ", prev, prevIndex)
            self.messageQueue[workflow].PutElement(msg, prevIndex)
            if self.messageQueue[workflow].IsReady():
                print("inputs are ready: ", workflow)
                res = self.messageQueue[workflow].GetList()
                data = str([ele.value.decode("utf-8") for ele in res])
                print(data)
                cmds = self.userCMD.split()
                cmds.append(data)
                userProcess = subprocess.Popen(cmds, stdout=subprocess.PIPE, text=True, universal_newlines=True)
                                        
                processOutput = str(userProcess.communicate()[0])
                print(processOutput, end='')

                nextContainers = output[containerName]
                
                
                msgID = msg.key.decode("utf-8").split("::")[2]
                key = workflow+"::"+containerName+"::"+str(msgID)
                for container in nextContainers:
                    print("sending to: ", container)
                    if container.startwiths("__") and container.endswith("__"):
                        self.producer.send(workflow+"-output",  key = key.encode("utf-8"),  value = processOutput.encode('utf-8'),partition=0, headers=[("Input", msg.headers[0][1]), ("Output", msg.headers[1][1]),('workflow', workflow.encode('utf-8'))])
                        continue

                    self.producer.send(container,  key = key.encode("utf-8"),  value = processOutput.encode('utf-8'),partition=0, headers=[("Input", msg.headers[0][1]), ("Output", msg.headers[1][1]),('workflow', workflow.encode('utf-8'))])
                # Otherwise call the user function and produce a next message

if __name__ == "__main__":
    se = streamExecutor()
    se.Exec()