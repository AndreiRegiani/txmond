FROM golang:1.20
WORKDIR /app
COPY . .
RUN go build -o txmond cmd/txmond.go
COPY .env .
EXPOSE 9901
CMD ["./txmond"]
