FROM golang:1.16 as builder
WORKDIR /go/src
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
COPY . /go/src
RUN go build -o linq .

FROM scratch as runner
COPY --from=builder /go/src/linq /app/linq
COPY ./test.yml /app/test.yml
