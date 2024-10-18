#!/bin/bash

# Exit script on any error
set -e

echo "Starting all backend services..."

# Function to start a Go service
start_go_service() {
    SERVICE_NAME=$1
    SERVICE_DIR=$2

    echo "Starting $SERVICE_NAME..."

    # Navigate to the service directory
    cd "$SERVICE_DIR"

    # Initialize Go module if not already initialized
    if [ ! -f go.mod ]; then
        echo "Initializing Go module for $SERVICE_NAME..."
        go mod init "$SERVICE_NAME"
    fi

    # Tidy up dependencies
    echo "Tidying up dependencies for $SERVICE_NAME..."
    go mod tidy

    # Run the service in the background and redirect output to log files
    echo "Running $SERVICE_NAME..."
    go run . > "../../logs/${SERVICE_NAME}.log" 2>&1 &

    # Return to the root directory
    cd - > /dev/null
}

# Create logs directory if it doesn't exist
mkdir -p logs

# Start backend services
start_go_service "payment_service" "backend/payment_service"
start_go_service "transport_service" "backend/transport_service"
start_go_service "transport_tracking_service" "backend/transport_tracking_service"
start_go_service "driver_location_service" "backend/driver_location_service"
start_go_service "transport_matcher_service" "backend/transport_matcher_service"

echo "All backend services have been started."

# Wait for all background jobs
wait
