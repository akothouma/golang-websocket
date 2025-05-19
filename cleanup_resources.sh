#!/bin/bash


IMAGE_NAME="forum"
CONTAINER_NAME="forum-container"


cleanup_project() {
    echo "Starting cleanup for ${IMAGE_NAME}..."

    # Stop and remove container (only if exists)
    if [ "$(docker ps -aq -f name= ${CONTAINER_NAME})" ]; then
        echo "Stopping and removing container: ${CONTAINER_NAME}"
        docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}
    else
        echo "No container named ${CONTAINER_NAME} found"
    fi

    # Remove image (only if exists)
    if [ "$(docker images -q ${IMAGE_NAME})" ]; then
        echo "Removing image: ${IMAGE_NAME}"
        docker rmi ${IMAGE_NAME}
    else
        echo "No image named ${IMAGE_NAME} found"
    fi

    # Clean up dangling build cache (safe)
    echo "Cleaning up dangling build cache..."
    docker builder prune -f

    echo "Cleanup completed for ${IMAGE_NAME}"
}

# Function to display usage
usage() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -c, --cleanup    Remove ONLY ${IMAGE_NAME} containers/images"
    echo "  -h, --help       Display this help message"
}

# Main execution
case "$1" in
    -c|--cleanup)
        cleanup_project
        ;;
    -h|--help)
        usage
        ;;
    *)
        echo "Invalid option: $1"
        usage
        exit 1
        ;;
esac