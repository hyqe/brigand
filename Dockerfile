FROM golang:bullseye as compile

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=compile /main /main

ENTRYPOINT ["/main"]