#!/bin/bash

# Set the image name and container name
IMAGE_NAME="forum"
CONTAINER_NAME="forum-container"

#  Build the Docker image
echo "Building Docker image: $IMAGE_NAME"
if ! docker build -t $IMAGE_NAME . ; then
    echo "Error: Docker build failed."
    exit 1
fi

# Check if the container already exists
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "Container '$CONTAINER_NAME' already exists."
    read -p "Do you want to remove it and create a new one? (y/n): " choice
    if [[ "$choice" == "y" ]]; then
        echo "Stopping and removing the existing container..."
        docker stop $CONTAINER_NAME
        docker rm $CONTAINER_NAME
    else
        echo "Please choose a different container name or remove the existing one manually."
        exit 1
    fi
fi

# Run the Docker container
echo "Running Docker container: $CONTAINER_NAME"
if ! docker run -d -p 8000:8000 --name $CONTAINER_NAME $IMAGE_NAME; then
    echo "Error: Docker run failed."
    exit 1
fi

# Inform the user that the application is running
echo "Your application is now running at http://localhost:8000                                                                                                 