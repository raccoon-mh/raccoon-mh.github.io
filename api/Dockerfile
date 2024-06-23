##############################################################
## Stage 1 - Go Build
##############################################################

# Using a specific version of golang based on alpine for building the application
FROM golang:1.22.4-alpine AS builder

RUN mkdir -p /src/raccoon-mh.github.io
WORKDIR /src/raccoon-mh.github.io

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

ADD . .
RUN go build

#############################################################
## Stage 2 - Deploy
##############################################################

# Using the latest Ubuntu image for the production stage
FROM ubuntu:latest

# Setting the working directory for the application 
WORKDIR /app

COPY --from=builder /src/raccoon-mh.github.io/api /app/api
RUN touch /app/.env

ENV ADDR 0.0.0.0
ENV PORT 3000

EXPOSE 3000

ENTRYPOINT [ "/app/api" ]
