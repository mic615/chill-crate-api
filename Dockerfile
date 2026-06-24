# official Golang image with Alpine Linux
FROM golang:alpine3.23 AS builder
# set working directory
WORKDIR /app
# copy dependencies
COPY go.mod go.sum ./
# download dependencies
RUN go mod download
COPY . .
# build the application
RUN go build -o chill-crate-api ./cmd/server
# build the migration tool
RUN go build -o migrate ./internal/migration/migrate.go
# create a minimal image for the final build
FROM alpine:3.23
# set working directory
WORKDIR /app
# copy the built application and migration tool from the builder stage
COPY --from=builder /app/chill-crate-api .
COPY --from=builder /app/migrate .
EXPOSE ${SERVER_PORT}
# run the migration tool and then start the application
CMD ["sh", "-c", "./migrate && ./chill-crate-api"]