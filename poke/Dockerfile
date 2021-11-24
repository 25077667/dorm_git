FROM golang:alpine as builder

WORKDIR /src
COPY ["go.mod", "go.sum", "main.go", "/src/"]
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main main.go

FROM alpine
LABEL MAINTAINER="zxc25077667@pm.me"
LABEL AUTHOR="25077667/SCC"
COPY --chown=1001:0 --from=builder /src/main /target/
WORKDIR /target
USER 1001

CMD ["./main"]