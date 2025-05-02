FROM golang:1.22

WORKDIR app/

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/dengovie

EXPOSE 8080

RUN useradd -ms /bin/bash dengovie-runner
USER dengovie-runner

CMD ["./dengovie"]