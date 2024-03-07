FROM golang:1.22.0-bookworm as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /liondb ./cmd/liondb

FROM debian:bookworm

COPY --from=build /liondb /usr/local/bin/liondb

CMD ["liondb"]
