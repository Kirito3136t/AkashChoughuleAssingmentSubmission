Libraries Used in this project

github.com/golang-jwt/jwt/v4 v4.5.2
github.com/golang-jwt/jwt/v5 v5.3.0
github.com/google/uuid v1.6.0
github.com/joho/godotenv v1.5.1
github.com/lib/pq v1.10.9
github.com/sirupsen/logrus v1.9.3

Install :

go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


Commands to setup database:

1. cd sql/schema
2. goose postgres postgres://<username>:<password>@<host>:<port>/assignment up
3. cd .. / cd..
4. sqlc generate

After using these commands just build the project and execute the apis