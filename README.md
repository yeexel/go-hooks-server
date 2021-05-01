# Go Hooks Server

![example workflow](https://github.com/yeexel/go-hooks-server/actions/workflows/actions.yml/badge.svg)

Simple hooks server written in Go with native `http` package.

### Installation

#### Local machine

Make sure you have Docker installed and then execute the following command:

`docker-compose up`

The server will automatically listen on port `9876`.

### Examples

The hooks server has 2 endpoints available.

#### POST /api/webhooks

This endpoint allows to register a webhook with `url`and `token` params provided.

Example cURL code snippet:

```
curl --location --request POST 'http://localhost:9876/api/webhooks' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url": "https://api.myserver.com/hooks",
    "token": "tok123"
}'
```

Example response:

```
{
  "id": "Qf2gRHt",
  "url": "https://api.myserver.com/hooks",
  "token": "tok123"
}
```

#### POST /api/webhooks/test

This endpoint is mainly used for testing hooks that were registered before.

`X-WebhookId` header is required.

Example cURL code snippet:

```
curl --location --request POST 'http://localhost:9876/api/webhooks/test' \
--header 'X-WebhookId: Qf2gRHt' \
--header 'Content-Type: application/json' \
--data-raw '{
    "payload": [{"user": "test"}]
}'
```

### Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/yeexel/go-hooks-server.svg)](https://pkg.go.dev/github.com/yeexel/go-hooks-server)

### Limitations

- The hooks server is using in-memory store which will be reset on server restart.
- Request timeout is missing for `/api/webhooks/test`endpoint meaning that requests to unknown/un-existing domains may take longer time to finish. This can be fixed by providing standalone HTTP client with `http.Transport` implementation.
- Rate-limiting is not supported at the moment.
