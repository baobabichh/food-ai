# Use an official C++ base image
FROM ubuntu:22.04

# Install dependencies
RUN prebuild.sh

# Set the working directory
WORKDIR /app

# Copy source files and CMakeLists.txt
COPY . .

# Create a build directory
RUN ./services/cpp/build_all.sh