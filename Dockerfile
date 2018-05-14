FROM golang:1.10

WORKDIR /go/src/github.com/arpradhan/usda-calorie-search
COPY . .

RUN go get -d -v \
    github.com/kelseyhightower/envconfig

RUN go build 
RUN go install -v ./app

CMD ["app"]