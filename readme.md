# Backend App - Simple Bank

> A working backend project base on Go, with good practice, using docker, postgresql

## Lauch Project

```bash
cd simplebank
docker compose up
```

Then it's possible to explore the API with [Postman](https://www.postman.com/) importing the following tests :
`postman/simple_bank.postman_collection.json`

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

## CRUD on Golang

We can use various tools with golang to work with databases.

Here are some things to consider to choose which one we will use.
![database tools](assets/golang-sql-tools.png)

[SQLC](https://sqlc.dev/) seems to be the more efficient.

We define CRUD methods in `db/query/` then generate the go functions with `make sqlc`

## Test the app

During the development, we can test the app with the help of some tools. For go, we can use [testify](https://github.com/stretchr/testify)

With the `-cover` option, we can see how much of our code is tested.

Also VSCode help us to see the coverage with the use of `run package test` who show us which part of our code are really tested.

It's possible to launch the test with the following command:

```bash
make test
```

## ACID (Atomicity, Consistency, Isolation, Durability)

![acid](assets/acid.png)

The transactions in a database <u>must respect</u> the ACID properties :

- Atomicity (A) : Either all operations complete successfully or the transaction fails and the db is unchanged.
- Consistency (C) : The db must be valid after the transaction. All constraints must be satisfied.
- Isolation (I) : Concurrent transactions must not affect each other.
- Durability (D) : Data written by a successful transaction must be recorded in persistent storage.

These process are implemented in db management systems like MySQL or PostgreSQL and there might be errors, timeout or deadlock who need to be managed in our go application.

### Comparation between MySQL and PostgreSQL :

| <u>MySQL:</u> | read uncommitted | read committed | reapeatable read | serialisable |
|---------------|------------------|----------------|------------------|--------------|
| dirty read           | ok        | x              |x                 | x            |
| non-reapeatable read | ok        | ok             | x                | x            |
| phantom read         | ok        | ok             | x                | x            |
|serialization anomaly | ok        | ok             | ok               | x            |

| <u>PostgreSQL:</u> | read uncommitted | read committed | reapeatable read | serialisable |
|--------------------|------------------|----------------|------------------|--------------|
| dirty read           | x              | x              |x                 | x            |
| non-reapeatable read | ok             | ok             | x                | x            |
| phantom read         | ok             | ok             | x                | x            |
|serialization anomaly | ok             | ok             | ok               | x            |

|                      |     MySQL          |     PostgreSQL          |
|----------------------|--------------------|-------------------------|
|                      | 4 isolation levels | 3 isolation levels      |
|                      | locking mechanism  | dependencies detections |
| deft isolation state | repeatable read    | read committed          |


## Choose a Framework or a HTTP router for our Golang project

There are a lot of projects that exist in Golang to help us to create a HTTP app.

![list of frameworks and http routers](assets/web_frameworks-http_routers.png)

In this project, we will choose [Gin](https://github.com/gin-gonic/gin) who is a really popular and fast web framework. It will also help us a lot to manage basics part of a HTTP application.

## A RESTful HTTP API

[<strong>Representational state transfer (REST)</strong>](https://en.wikipedia.org/wiki/Representational_state_transfer) is a software architectural style that describes a uniform interface between decoupled components in the Internet in a Client-Server architecture. REST defines four interface constraints:
- Identification of resources
- Manipulation of resources
- Self-descriptive messages and
- hypermedia as the engine of application state

Semantics of HTTP methods:
- GET:  Get a representation of the target resource’s state.
- POST: Let the target resource process the representation enclosed in the request.
- PUT:  Create or replace the state of the target resource with the state defined by the representation enclosed in the request.
- PATCH: Partially update resource’s state.
- DELETE: Delete the target resource’s state.
- OPTIONS: Advertising the available methods.

## Load config from file & environment variables

When we work with a docker or various environment, it is a good practice to separate some information from the rest of the code and make our system works with `ENV VARS` for a better organization.

For that purpose, we can use a tool like [viper](https://github.com/spf13/viper) to help us.

Now we can override our configuration like that :

```bash
$ SERVER_ADDRESS=0.0.0.0:8081 make server
...
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8081

# or
$ make server -e SERVER_ADDRESS=0.0.0.0:8081
...
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8081
```

## Mock DB for testing HTTP API

A mock DB is a fake database with the same definition of our real DB, in order to make all test on it without modifying the working environment.

There are various way to implement that:

- ### Use fake DB: Memory

> Implement a fake version of BD: store data in memory:
> - (+) simple approach
> - (+) easy to implement
> - (-) require lot of code for testing which is time consuming for development and maintenance later

```go
type Store interface {
    GetAccount(id int64) (Account, error)
}

type MemStore struct {
    data map[int64]Account
}

func (store *MemStore) GetAccount(id int64) (Account, error) {
    return store.data[id], nil
}
```

- ### Use DB Stubs: GoMock
> Generate and build stubs that returns hard-coded values
>
> https://github.com/golang/mock

```go
func TestGetAccountAPI(t *testing.T){
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    store := mockdb.NewMockStore(ctrl)
    account := randomAccount()

    // Build stub: expect GetAccount() to be called
    // exactly 1 time with this input account.ID
    // and return this account object as output
    store.EXPECT().GetAccount(gomock.Eq(account.ID)).Times(1).Return(account)

    // TODO: use this mock store to test your API
}
```

### Generate mock DB

We can generate the mock with the following command:

```bash
make mock
```
