# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Copy the images from the host into the container
COPY images /app/images

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8181

# Run the Go application
CMD ["./main"]

