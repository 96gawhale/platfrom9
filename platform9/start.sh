#!/bin/bash

# Function to start a server and save its PID
start_server() {
    local server_dir=$1
    local server_name=$2
    local port=$3
    
    # Navigate to the server directory
    cd "$server_dir" || exit
    
    # Start the server in the background without generating log files
    nohup go run main.go > /dev/null 2>&1 &
    
    # Print server started message
    echo "$server_name server started on http://localhost:$port"
    
    # Return to the original directory
    cd - > /dev/null || exit
}

# Start each server
start_server "." "Main" 8080
start_server "cal" "Calculator" 8081
start_server "todo" "Todo" 8082

