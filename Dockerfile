FROM node as docs
WORKDIR /docs
COPY . .
RUN scripts/compile_openapi.sh


FROM golang:bullseye as compile
WORKDIR /app
COPY . .
COPY --from=docs /docs/internal/handlers/docs.html internal/handlers/docs.html
RUN go vet ./...
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=compile /main /main
ENTRYPOINT ["/main"]