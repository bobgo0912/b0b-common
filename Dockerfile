FROM alpine:3.14

ARG B0B_ENV
WORKDIR /app
ENV B0B_ENV=${B0B_ENV}
COPY app ./
COPY config ./config

EXPOSE 8888

CMD [ "./app" ]