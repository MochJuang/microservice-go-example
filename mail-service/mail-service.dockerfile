FROM alpine:latest

RUN mkdir /app

COPY mailApp  /app

RUN mkdir /templates
COPY ./templates/ /templates

CMD ["/app/mailApp"]