FROM golang:1.21-alpine
 
# Create a directory for the app
RUN mkdir /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
 
# Seting working directory
WORKDIR /app

COPY . .

# Building the Go application
RUN go build -v -o blog .

# Run the blog executable
CMD [ "/app/blog" ]
