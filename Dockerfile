FROM golang:alpine

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o scraper

CMD ["/app/scraper"]

