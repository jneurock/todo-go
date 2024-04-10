FROM golang:1.22

WORKDIR /usr/src/todo

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir /usr/local/bin/todo
RUN go build -v -o /usr/local/bin/todo ./...

ENTRYPOINT ["/usr/local/bin/todo/web"]
