FROM golang:latest

RUN mkdir /go/src/go-todo

ADD main.go /go/src/go-todo
ADD handler /go/src/go-todo/handler
ADD model /go/src/go-todo/model

WORKDIR /go/src/go-todo
EXPOSE 1192

ENV DEBUG TRUE

RUN go get 	github.com/gin-gonic/contrib/sessions
RUN go get 	github.com/gin-gonic/gin
RUN go get 	github.com/go-sql-driver/mysql
RUN go get 	github.com/jinzhu/gorm
RUN go get 	github.com/jinzhu/gorm/dialects/mysql

CMD ["go", "run", "main.go"]