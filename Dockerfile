# docker run --net=host -it --rm server:v1

FROM golang:1.23.1-alpine3.19 as base
FROM base as dev
WORKDIR /app
COPY . .
EXPOSE 8082
RUN apk update
RUN apk add make 
CMD ["make",  "run"]