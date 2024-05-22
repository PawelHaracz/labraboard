FROM golang:1.22.0-alpine3.19 as build

ENV GO111MODULE=on

WORKDIR /app/build

# Cache go.mod for downloading dependecies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# App binary without CGO_ENABLED
ENV CGO_ENABLED=0 GOOS=linux

# Find and build all main.go files
RUN for file in $(find /app/build/cmd -name main.go); do \
        go build -o /app/build/bin/$(dirname $file | xargs basename) $file; \
    done

FROM alpine:edge
WORKDIR /app

COPY entrypoint.sh entrypoint.sh

#COPY --from=build /app/build/cmd/ ./
COPY --from=build /app/build/bin/ /app/

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set exec permision
RUN find /app -type f -exec chmod +x {} \;

# Run binary as non-root
RUN addgroup --system runner && adduser --system --no-create-home --disabled-password runner && adduser runner runner
USER runner
EXPOSE 8080

ENTRYPOINT ["sh", "./entrypoint.sh"]
CMD ["./api" ]
#ENTRYPOINT ["./cmd/api/main"]