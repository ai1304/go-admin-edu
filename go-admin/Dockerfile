FROM docker.1ms.run/library/golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache gcc g++ musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/go-admin .

FROM docker.1ms.run/library/alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=build /out/go-admin /main
COPY ./config /config
EXPOSE 8000
RUN chmod +x /main
CMD ["/main", "server", "-c", "/config/settings.yml"]
