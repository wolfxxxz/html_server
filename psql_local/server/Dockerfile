FROM golang:1.20-alpine


COPY . /server
WORKDIR /server

RUN go mod tidy
RUN go build -o main cmd/server/main.go
CMD [ "./main" ]


