FROM alpine:3.14

WORKDIR /app
COPY bin/app /app/
COPY config/config.env /config/config.env
RUN chmod +x /app/app
EXPOSE 8080
CMD ["./app"]
