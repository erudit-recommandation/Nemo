FROM golang:1.18

RUN mkdir /app

ADD . /app/

WORKDIR /app

EXPOSE 8087

CMD ["make", "run"]