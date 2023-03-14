FROM golang:latest
# Start from a Golang image

# Install git, required for fetching the dependencies
RUN apt-get update


WORKDIR /home/dev1/GolandProjects/atm-machine
COPY . .
RUN pwd
RUN ls

RUN cd services/discovery
RUN pwd

RUN ls

RUN go build -o discovery ./services/discovery

EXPOSE 8091
CMD ["./discovery"]

