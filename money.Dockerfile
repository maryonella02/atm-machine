FROM golang:latest
# Start from a Golang image

WORKDIR /home/dev1/GolandProjects/atm-machine

COPY . .

# Install git, required for fetching the dependencies
RUN apt-get update

# Install the dependencies
RUN go get github.com/prometheus/client_golang/prometheus@latest


RUN cd services/money_operations
RUN go build -o money-operations ./services/money_operations

EXPOSE 8092 8093
CMD ["./money-operations"]

