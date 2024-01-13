FROM golang:1.21.5
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-playground

EXPOSE 8080

CMD ["/go-playground"]