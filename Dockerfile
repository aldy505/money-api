FROM golang:1.17.2-buster

WORKDIR /app

ENV PORT=3000

COPY . .

RUN apt-get update && \
  apt-get upgrade -y  && \
  apt-get install -y sqlite3

RUN go mod download

RUN go build main.go

EXPOSE ${PORT}

CMD [ "./main" ]