# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

# Download required modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source and gp2midi tool
COPY GuitarProToMidi ./
COPY *.go ./

# Build
RUN go build -o /gp2midi-web

# Application port
EXPOSE 8229

# Run
CMD [ "/gp2midi-web" ]
