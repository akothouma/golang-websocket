FROM golang:1.23-alpine AS builder

WORKDIR  /app

COPY . .

RUN go build -o main .

LABEL maintainer="lornaakothpb@gmail.com"

FROM debian:bookwarm-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY ui /app/ui  


EXPOSE 8000

CMD [ ./main ]