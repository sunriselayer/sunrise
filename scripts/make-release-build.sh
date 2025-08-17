#!/bin/bash
set -eux

# Create the output directory if it doesn't exist
mkdir -p ./docker-build

# Remove the old binary if it exists
sudo rm -f ./docker-build/sunrised

# Build the Docker image
docker build . -f Dockerfile -t sunrised-dev

# Run the container and copy the built binary
docker run --rm -v $PWD/docker-build:/root sunrised-dev sh -c "cp /sunrise/build/sunrised /root/"

# Print the md5sum of the binary
md5sum ./docker-build/sunrised

# Check the version of the binary
./docker-build/sunrised version
