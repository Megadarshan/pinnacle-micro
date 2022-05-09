FROM golang:1.17-alpine as builder
# RUN mkdir /service
WORKDIR /service
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# RUN go build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o managetoken
# CMD [ "./redis" ]

# image.
FROM alpine:latest
# Security related package, good to have.
RUN apk --no-cache add ca-certificates
# COPY — from=build /redis /redis
RUN mkdir /app
WORKDIR /app
COPY --from=builder /service/managetoken .
# ENTRYPOINT [ "/redis" ]
ENTRYPOINT [ "./managetoken" ]