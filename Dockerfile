#              MULTI STAGE BUILD FOR COMPILED LANGUAGES

# STEP 1 BUILD STAGE

# INSTRUCTIONS args for base image
FROM golang:1.24-alpine AS builder

#TO imporve speed of image
RUN apk add curl
RUN apk add --no-cache git

#setting up working directory
WORKDIR /Downloads

#copying dependencies
COPY go.mod go.sum ./
RUN go mod download

#copying all code
COPY . .

#build
RUN go build -o relay main.go

#RUN STAGE 2
FROM alpine:latest

#working dir
WORKDIR /Downloads

#
COPY --from=builder /Downloads/relay .

# Exposing the port (same as my PORT env)
EXPOSE 443

# Run the relay server
CMD ["./relay"]

