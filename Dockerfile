FROM golang:1.16
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux GO111MODUL=on go build -o tf-task-runner .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /src/tf-task-runner .
ENTRYPOINT ["./tf-task-runner"]  