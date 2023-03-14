FROM golang:latest
# Start from a Golang image

# Install git, required for fetching the dependencies
RUN apt-get update


WORKDIR /home/dev1/GolandProjects/atm-machine
COPY . .
RUN pwd
RUN ls

RUN cd services/gateway
RUN go build -o gateway ./services/gateway

EXPOSE 8090
CMD ["./gateway"]

