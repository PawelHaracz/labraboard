FROM golang:1.22.0-alpine3.19 as build

ENV GO111MODULE=on

WORKDIR /app/build

# Cache go.mod for downloading dependecies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# App binary without CGO_ENABLED
RUN CGO_ENABLED=0 GOOS=linux go build /app/build/cmd/main.go

FROM alpine:edge
WORKDIR /app

COPY --from=build /app/build/main ./
# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set exec permision
RUN chmod +x ./main

# Run binary as non-root
RUN addgroup --system runner && adduser --system --no-create-home --disabled-password runner && adduser runner runner
USER runner
EXPOSE 8080
ENTRYPOINT ["./main"]