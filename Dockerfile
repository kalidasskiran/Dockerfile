FROM golang:latest
# Add a work directory


WORKDIR /app

COPY . .

RUN go build main.go

CMD ["/app/main"]