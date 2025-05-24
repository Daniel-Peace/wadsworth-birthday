## Manual DB Connection
To connect to the db through the terminal you can use:
```
mongosh mongodb://localhost:8000
```

## Switch db
To switch to the which database you are viewing use:
```
use <db>
```
Where db is the database you want to view. (In this case use `wadsworth-birthday` for the db)

## Manual Collection query
To find a collection manually you can use the following:
```
db.<collection>.find()
```
where `<collection>` is your collection's name (In this case use `birthdays` for the collection)

## Starting the Docker container
If the docker container is ever stopped, you can run the following to restart it:
```
sudo docker start wadsworth-db
```
