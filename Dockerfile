# Start from a Golang base image
FROM golang:1.22.2

# Set the current working directory inside the container
WORKDIR /go/src/app

# Copy everything from the current directory to the WORKDIR inside the container
COPY . .

# Build the Go app
RUN make build

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["make", "run"]