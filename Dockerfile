FROM golang:1.10.2
LABEL maintainer="Cristian C. Wolfram <cristiancw@gmail.com>"

ENV GOBIN $GOPATH/bin

WORKDIR ./
COPY *.go ./
COPY info/*.go ./info/

RUN go get -d -v github.com/gorilla/mux github.com/gocql/gocql
RUN go install -v ./myinfos.go

EXPOSE 8080

ENTRYPOINT [ "myinfos" , "-port=8080", "-hostdb=local.myinfos", "-portdb=9042"]
