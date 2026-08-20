FROM golang:1.22-bookworm
ENV GOPROXY=https://goproxy.cn,direct
ENV GOTOOLCHAIN=local
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /usr/local/bin/bathypulse .
EXPOSE 8097
ENTRYPOINT ["/usr/local/bin/bathypulse"]
