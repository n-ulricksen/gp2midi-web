# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

# Download required modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Install packages to run gp2web program
RUN apt -y update
RUN apt -y install libicu-dev

# Copy source and gp2midi tool
COPY *.go ./
COPY GuitarProToMidi ./

# Build
RUN go build -o /gp2midi-web

# Expose port (only used in dev)
EXPOSE 8229

# Run
CMD [ "/gp2midi-web" ]
