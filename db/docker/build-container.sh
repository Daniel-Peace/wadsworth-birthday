#!/bin/bash
echo "stopping running container if it exists..."
sudo docker stop wadsworth-db

echo "deleting container if it exists..."
sudo docker rm wadsworth-db

echo "copying piblic ssh key into local directory"
sudo cp /home/pacodataco/.ssh/id_rsa.pub ./

echo "building container..."
sudo docker build -t wadsworth-db-image .

echo "removing ssh key from local directory"
sudo rm ./id_rsa.pub

echo "done!"
