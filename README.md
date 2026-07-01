Steps to run:
go mod init golang-postgres-crud

go get -u github.com/gin-gonic/gin
go get -u github.com/jackc/pgx/v5/pgxpool
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files


swag init

go mod tidy
