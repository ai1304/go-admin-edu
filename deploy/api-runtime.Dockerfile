FROM docker.1ms.run/library/alpine:3.21

WORKDIR /app
COPY go-admin /main
EXPOSE 8000

RUN chmod +x /main
CMD ["/main", "server", "-c", "/config/settings.yml"]
