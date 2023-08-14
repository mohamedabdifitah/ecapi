FROM golang:1.19.3-alpine3.16
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 80:80
CMD ["/app/main"]