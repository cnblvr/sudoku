FROM golang:latest

WORKDIR /usr/src/sudoku

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/sudoku ./sudoku/cmd/main.go

CMD ["sudoku"]