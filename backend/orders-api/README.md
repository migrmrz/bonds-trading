# orders-api
API that lets users create, insert and update limit orders.

## Endpoints
### UpsertOrders
This lets users create a limit order (if it doesn't exist) or update an existing one (by its `bond_id`) like price or status. For an update, all fields need to be sent, along with those that need update.

**Note:** Both endpoints require basic authentication

`PUT /book/v1/orders`
#### Example HTTP Request (insert)
```
PUT /book/v1/orders HTTP/1.1
Host: localhost:8000
Content-Type: application/json
Authorization: Basic YWRtaW46bWFuYWdlci4x
Content-Length: 172

{
    "bond_id": "20240111131334",
    "action": "sell",
    "quantity": 45,
    "price": 10000,
    "status": "active",
    "user": "admin",
    "created_at": 1704970813
}
```

##### cURL
```curl
curl --location --request PUT 'http://localhost:8000/book/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic YWRtaW46bWFuYWdlci4x' \
--data '{
    "bond_id": "20240111131334",
    "action": "sell",
    "quantity": 45,
    "price": 10000,
    "status": "active",
    "user": "admin",
    "created_at": 1704970813
}'
```

#### Example Response
```json
{
    "data": "bond:20240111131334"
}
```
---
#### Example HTTP Request (update)
```
PUT /book/v1/orders HTTP/1.1
Host: localhost:8000
Content-Type: application/json
Authorization: Basic YWRtaW46bWFuYWdlci4x
Content-Length: 172

{
    "bond_id": "20240111131334",
    "action": "sell",
    "quantity": 45,
    "price": 10000,
    "status": "cancelled",
    "user": "admin",
    "created_at": 1704970813
}
```

##### cURL
```curl
curl --location --request PUT 'http://localhost:8000/book/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic YWRtaW46bWFuYWdlci4x' \
--data '{
    "bond_id": "20240111131334",
    "action": "sell",
    "quantity": 45,
    "price": 10000,
    "status": "active",
    "user": "admin",
    "created_at": 1704970813
}'
```
#### Example Response
```json
{
    "data": "bond:20240111131334"
}
```

### GetOrders
GetOrders receives either the `user` parameter to filter by user or the `status` parameter to filter by status (active).

`GET /book/v1/orders`
#### Example HTTP Request (filter by user)
```
GET /book/v1/orders?user=admin HTTP/1.1
Host: localhost:8000
Authorization: Basic YWRtaW46bWFuYWdlci4x
```

##### cURL
```curl
curl --location 'http://localhost:8000/book/v1/orders?user=admin' \
--header 'Authorization: Basic YWRtaW46bWFuYWdlci4x'
```
#### Example HTTP Request (filter by status)
```
GET /book/v1/orders?status=active HTTP/1.1
Host: localhost:8000
Authorization: Basic YWRtaW46bWFuYWdlci4x
```

##### cURL
```curl
curl --location 'http://localhost:8000/book/v1/orders?status=active' \
--header 'Authorization: Basic YWRtaW46bWFuYWdlci4x'
```

## How to run the service

The most useful make targets for working locally are:

* `make build`: Builds the service.
* `make run`: Starts the service locally running on port `8001`.
* `make clean`: Clean temporary files.

## Dependencies
This project has dependencies on:
* go (`1.20.12`)
* PostgreSQL (`16.1`)
* Redis (`7.2.3`)
* NATS (`1.31.0`)

### PostgreSQL
1. Start a local instance or a docker container with a [postgres docker image](https://hub.docker.com/_/postgres) using default port `5432`
2. Create `ordersdb` database
3. Create USERS table
   ```sql
   CREATE TABLE USERS (user_id serial PRIMARY KEY,
     username varchar(50) UNIQUE NOT NULL,
     hashed_password varchar(100) NOT NULL,
     email varchar(100) UNIQUE NOT NULL,
     created_at timestamp NOT NULL,
     last_login timestamp);
   ```
4. Insert a user
   ```sql
   INSERT INTO USERS(username, hashed_password, email, created_at)
   VALUES ('admin','$2a$10$2RN.gglXz6ZeHtNmdoUgDe4h7rUGNlU8qsu70N7K9U4sBSKtbZzaO','admin@trading.com.mx',CURRENT_TIMESTAMP);
   ```

### Redis
Start a local instance or a docker container with a [redis docker image](https://hub.docker.com/_/redis) using port `6379`

### NATS
Create and start a local server, or you can also use [docker](https://hub.docker.com/_/nats) for this with default port `4222`
