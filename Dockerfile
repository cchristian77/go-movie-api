FROM alpine:latest

RUN mkdir /app

WORKDIR /app

RUN mkdir /configs

COPY /configs/env.json /app/configs/env.json

COPY go-movie-api-build /app

CMD [ "/app/go-movie-api-build" ]