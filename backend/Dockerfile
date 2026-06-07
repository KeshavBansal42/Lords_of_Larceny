FROM golang:1.26.4

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

EXPOSE 3000

CMD [ "./server" ]