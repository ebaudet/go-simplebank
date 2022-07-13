# Backend App - Simple Bank

> A working backend project base on Go, with good practice, using docker, postgresql

## Diagram

Database diagram is done on [dbdiagram.io](https://dbdiagram.io/d/62cd1dc8cc1bc14cc59dd2c5)

## Docker

postgres images : postgres:14-alpine
https://hub.docker.com/_/postgres


```bash
# pull image
docker pull postgres:14-alpine

# start a container
docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

# run into it
docker exec -it postgres14 psql -U root

# check logs
docker logs postgres14

# stop container
docker stop postgres14
# start container
docker start postgres14
```

We can manage our postgres database with tools like [TablePlus](https://tableplus.com/).

## Database migration

We can use this [migrate tool](https://github.com/golang-migrate/migrate/) to manage our migration, version of database.

```bash
migrate create -ext sql -dir db/migration -seq init_schema

migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

then fill the migration up with our schema and down with `drop table`.

## Makefile instructions

```bash
# create and launch container
make postgres
# create database
make createdb
# push the migrations
make migrateup
```

