FROM golang:1.17-alpine

RUN mkdir /docker_back
ADD . /docker_back
WORKDIR /docker_back

RUN go build -o main_server .

EXPOSE 8080 8080
CMD ["/docker_back/main_server"]