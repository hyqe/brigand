FROM node as docs
WORKDIR /app
COPY . .
RUN scripts/compile_openapi.sh


FROM golang:bullseye as compile
WORKDIR /app
COPY . .
RUN go vet ./...
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=docs /app/openapi.html /openapi.html
COPY --from=compile /main /main
ENTRYPOINT ["/main"]