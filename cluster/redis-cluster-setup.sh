#!/bin/bash

# Initialize Redis cluster
echo "Initializing Redis cluster..."
IP_ADDRESS="$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}:6379' $(docker ps -q))"
docker run --rm --network hello-redis_redis-cluster -it redis redis-cli --cluster create $IP_ADDRESS --cluster-replicas 1

echo "Redis cluster initialized."
