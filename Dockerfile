FROM golang:1.19-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY rsscombine.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/rsscombine ./rsscombine.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=build /out/rsscombine /usr/local/bin/rsscombine

WORKDIR /app
ENTRYPOINT ["/usr/local/bin/rsscombine"]
