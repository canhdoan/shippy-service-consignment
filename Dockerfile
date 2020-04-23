# Dockerfile use multi-stage builds and run the image container
# https://docs.docker.com/develop/develop-images/multistage-build/

# Build image
FROM golang:alpine AS builder

RUN apk update \ 
    && apk upgrade \
    && apk add --no-cache git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shippy-service-consignment .

# Run containter
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/shippy-service-consignment .

CMD ["./shippy-service-consignment"]
