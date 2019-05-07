# use scratch for K8S production
ARG BASEIMAGE=scratch

FROM golang:1.12 as builder
WORKDIR /go/src/github.com/frnksgr/fibo
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install ./...

# NOTE: cf requires more than scratch
# while K8S is fine with it.
# build image for cf with 
# docker build -t frnksgr/fibo-cf --build-arg BASEIMAGE=alpine:3.9 .

FROM $BASEIMAGE
COPY --from=builder /go/bin/fibo /fibo
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/fibo"]
