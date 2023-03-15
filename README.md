# ATM Machine


The ATM Machine is a simple distributed system that simulates an ATM machine. It is comprised of four microservices:

    Discovery Service: A service registry that keeps track of all running microservices.
    Money Operations Service: A service that handles withdrawal, and balance inquiry requests.
    Gateway Service: A service that serves as a gateway for incoming client requests.
    Client Service: A client that sends requests to the Gateway Service.
    

### Getting Started

To get started with the ATM Machine, you will need to have the following tools installed:

    Go (v1.16 or later)
    Docker (v20.10 or later)
    Docker Compose (v1.29 or later)

To run the ATM Machine, follow these steps:

  1.Clone the repository to your local machine.
  
  2.Open a terminal and navigate to the project directory.
  
  3. Run the following command to build the Docker images.
      
      >docker-compose build
        
  4. Run the following command to start the Docker containers.

      >docker-compose up
       

This will start the Discovery Service, Money Operations Service, Gateway Service, and Client Service, Prometheus and Grafana.


### Monitoring

The ATM Machine uses Prometheus for monitoring. Grafana is used to visualize the metrics collected by Prometheus.

To access the Grafana dashboard, open a web browser and navigate to http://localhost:3000. Log in with the username admin and the password admin. 
