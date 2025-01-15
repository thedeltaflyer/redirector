FROM golang:1.23 AS builder

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . /app

# Run all go commands to build and test.
#    * Tidy up
#    * Run any generation that needs to run
#    * Run any unit tests included (with race detection)
#    * Build as a static executable
RUN go mod tidy \
    && go generate \
    && go test ./... -v -race -cover \
    && CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags "-s -w -extldflags \"-static\"" -buildvcs=false -o redirector


FROM scratch

WORKDIR /app

# Expose the port
EXPOSE 8080

# Copy the main Binary to the root folder
COPY --from=builder /app/redirector /app

# Set the entrypoint to run the application,
# if your application uses flags, make sure to
# add them here.
ENTRYPOINT ["./redirector"]
