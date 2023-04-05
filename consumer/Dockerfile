FROM golang:1.20-alpine as base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
ENV GIN_MODE=release
RUN go build -o /main

FROM alpine:3.17.2 as runtime
COPY --from=base /main /
RUN mkdir data
CMD [ "/main" ]

