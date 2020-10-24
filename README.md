# atlant.io

## Description

Go grpc service with mongo persistence running in a docker compose environment with Kong based load balancing.

1. Build the image:

```sh
make build
```

2. Start the containers (2 by default)

```sh
make compose
```

3. Provision the load balancer

```sh
make provision
```

4. Open your grpc client and point it to `127.0.0.1:9080` (no TLS) and import the `product.proto` descriptor file


## Testing
To load a demo list of products issues the follwing request:

```json
{
  "url": "http://csvhosting:8080/mockdata.csv"
}
```

To list by name:

```json
{
  "sortBy": "Name",
  "sortDesc": false,
  "page": 3,
  "limit": 10
}
```

or by update count:


To list by name:

```json
{
  "sortBy": "UpdateCount",
  "sortDesc": true,
  "page": 0,
  "limit": 3
}
```