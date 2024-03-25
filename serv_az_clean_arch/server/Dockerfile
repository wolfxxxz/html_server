FROM golang:1.20-alpine


COPY . /server
WORKDIR /server

RUN go mod tidy
RUN go build -o server cmd/server/main.go
CMD [ "./server" ]


