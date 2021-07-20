FROM golang:alpine

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
COPY --chown=root permissions/meetup-crawler-store-b25be2c787ec.json ./creds.json
RUN go build -v -o scraper

CMD ["/app/scraper"]

