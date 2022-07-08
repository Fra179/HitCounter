FROM golang:alpine AS builder
RUN mkdir "/project"
WORKDIR "/project"
RUN apk add build-base
COPY . .
RUN go build -o main .

FROM alpine:latest
COPY --from=builder "/project/main" .
EXPOSE 8080
RUN chmod +x main
HEALTHCHECK --interval=30s --timeout=5s --retries=10 CMD curl --fail http://localhost:$PORT/status || exit 1
CMD ["/main"]