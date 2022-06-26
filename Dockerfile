#builder
FROM golang:alpine as builder
WORKDIR /home
COPY . .
RUN go build -o storage-api app/main.go

#final image
FROM alpine
RUN apk add tzdata
COPY --from=builder /home/storage-api .
EXPOSE 6004
CMD ./storage-api