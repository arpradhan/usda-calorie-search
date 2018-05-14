FROM golang:1.10

WORKDIR /go/src/app
COPY . .

RUN go get -d -v \
    github.com/arpradhan/usda-calorie-search \
    github.com/kelseyhightower/envconfig

RUN go install -v ./app

CMD ["app"]