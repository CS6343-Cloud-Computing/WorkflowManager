#! /bin/bash
echo "Killing docker container"
sudo docker kill $(docker ps -q)
echo "Running Docker Container Prune"
sudo docker container prune -f
echo "Removing the images"
sudo docker rmi -f $(docker images -aq)
====== The End ======
====== The End ======