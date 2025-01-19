#!/bin/bash

# Build the Docker image
docker-compose build

# Run the Docker containers
docker-compose up -d

# Print the logs
docker-compose logs -f
