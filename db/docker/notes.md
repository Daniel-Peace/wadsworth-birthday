## Manual DB Connection
To connect to the db through the terminal you can use:
```
mongosh mongodb://localhost:8000
```

## Manual Collection query
To find a collection manually you can use the following:
```
db.<collection>.find()
```
where `<collection>` is your collection's name

## Starting the Docker container
If the docker container is ever stopped, you can run the following to restart it:
```
sudo docker start wadsworth-db
```
