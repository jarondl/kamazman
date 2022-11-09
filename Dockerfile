
# This is just to build the app.
FROM golang:1.19-bullseye as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY kamazman.go ./

RUN go build -v -o kamazman



FROM  debian:bullseye-slim
COPY --from=builder /app/kamazman /app/kamazman
COPY static/ /app/static/
COPY arrivals.db /app/

WORKDIR /app
EXPOSE 8080
CMD ["/app/kamazman"]
