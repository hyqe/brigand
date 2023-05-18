# Brigand

A file service


Run Brigand locally

```sh
MONGO="mongodb://localhost:27017" \
    go run main.go
```

Run locally using Docker compose. suitable for testing only.

```sh
docker compose -f ./scripts/docker-compose-local.yaml --project-directory . up  --build

# ctrl+c to stop
docker compose -f scripts/docker-compose-local.yaml  --project-directory . down
```