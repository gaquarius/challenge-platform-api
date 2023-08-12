# FROM golang:latest

# ENV GO111MODULE=on
# ENV PORT=8080
# WORKDIR /app
# COPY go.mod /app
# COPY go.sum /app

# RUN go mod download
# RUN go install -mod=mod github.com/githubnemo/CompileDaemon
# COPY . /app
# ENTRYPOINT CompileDaemon --build="go build -o main" --command=./main

FROM golang:latest

RUN mkdir /golang

RUN go install github.com/cosmtrek/air@latest

ADD . /golang/

RUN go install github.com/cosmtrek/air@latest

WORKDIR /golang

RUN go mod download

CMD ["air", "-c", ".air.toml"]