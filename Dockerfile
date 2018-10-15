FROM golang:1.11-stretch

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 


RUN go get -u	"github.com/Pallinder/sillyname-go"
RUN go get -u "github.com/mongodb/mongo-go-driver/bson"
RUN go get -u "github.com/mongodb/mongo-go-driver/mongo"

RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o server main.go

EXPOSE 5000

CMD ["./server"]


