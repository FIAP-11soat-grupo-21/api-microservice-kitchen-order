FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY . .

RUN go mod download
RUN go build -buildvcs=false -o main.exe .

FROM alpine:latest AS runtime

COPY --from=builder src/main.exe .

ENTRYPOINT [ "./main.exe" ]
