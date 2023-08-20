# go-example

golang [gin,gorm]

### start db & app

```start db & app
docker compose up -d
go run main.go
```

### stop db

```stop db
docker compose down
```

```
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/joho/godotenv
```
