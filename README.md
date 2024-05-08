# Eniqilo Store Back End

Requirement: [Project Sprint Inventory Management](https://openidea-projectsprint.notion.site/EniQilo-Store-93d69f62951c4c8aaf91e6c090127886?pvs=4)

Make new migration script

```migrate create -ext sql -dir migration -seq init```

Migrate database

up : ```migrate -database "postgres://postgres:password@localhost:5432/eniqilo-store?sslmode=disable" -path db/migrations up```

down : ```migrate -database "postgres://postgres:password@localhost:5432/eniqilo-store?sslmode=disable" -path db/migrations down```


## Run & Build Eniqilo Store

Run for debugging

```
make run
```

Build app

```
make build
```



