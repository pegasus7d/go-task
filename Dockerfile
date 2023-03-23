FROM golang:1.19-alpine
COPY . /app/
WORKDIR /app
RUN source .env
RUN go get github.com/jinzhu/gorm github.com/lib/pq github.com/gorilla/mux github.com/jinzhu/inflection
CMD go run main.go