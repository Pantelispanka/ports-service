# # specify the base image to  be used for the application, alpine or ubuntu
# FROM golang:1.17-alpine

# # create a working directory inside the image
# WORKDIR /app

# # copy Go modules and dependencies to image
# COPY go.mod ./

# # download Go modules and dependencies
# RUN go mod download

# # copy directory files i.e all files ending with .go
# COPY *.go ./

# # compile application
# RUN go build -o /goports

# # tells Docker that the container listens on specified network ports at runtime
# EXPOSE 8080

# # command to be used to execute when the image is used to start a container
# CMD [ "/goports" ]

FROM golang:alpine AS builder

WORKDIR /app

COPY ../go.mod go.sum ./

RUN go mod download

COPY .. .

RUN go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app .

CMD ["./app"]