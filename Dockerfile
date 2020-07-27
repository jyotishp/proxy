FROM golang:alpine as builder

RUN mkdir -p /app
WORKDIR /app
COPY . /app/
RUN go mod download
RUN go build -o proxy main.go


FROM scratch as proxy
COPY --from=builder /app/proxy /bin/proxy
EXPOSE 8080
ENTRYPOINT [ "proxy" ]
