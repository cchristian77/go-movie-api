FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY env.json /app

COPY go-movie-api-build /app

CMD [ "/app/go-movie-api-build" ]