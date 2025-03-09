#!/bin/bash
#
# scale_test.sh
#
# This script performs scale testing for the LB using all three policies:
#   - pick_first
#   - round_robin
#   - least_load
#
# For each policy, it:
#   1) Starts the LB (from server/lb.go), spawns multiple servers (from server/server.go),
#      and launches multiple clients (from client/client.go).
#   2) Runs for a specified duration (TEST_DURATION seconds).
#   3) Automatically terminates all processes.
#   4) Analyzes logs to compute performance metrics:
#         - Total responses (from client.log)
#         - Throughput (responses per second)
#         - LB assignment distribution (using both "Assigned server:" and policy-specific log messages)
#   5) Appends the results to a final summary log file.
#
# Adjust the configuration variables below as needed.

# -----------------------------
# CONFIGURATION
# -----------------------------
# List of policies to test:
LB_POLICIES=("pick_first" "round_robin" "least_load")
LB_PORT=9000
NUM_SERVERS=5                # Number of servers to spawn
START_PORT=9100              # Starting port for servers (9100, 9101, …)
NUM_CLIENTS=10               # Number of client processes to launch
CLIENT_REQUESTS=20           # Each client sends this many requests
CLIENT_CONCURRENCY=2         # Number of concurrent goroutines per client
TEST_DURATION=60             # Duration (in seconds) to run each test
FINAL_LOG="final_performance.log"

# Remove any previous final log file
rm -f "$FINAL_LOG"

# -----------------------------
# MAIN LOOP: Test Each Policy
# -----------------------------
for POLICY in "${LB_POLICIES[@]}"; do
  echo "========================================="
  echo "Running test for policy: $POLICY"
  echo "========================================="

  # Remove old logs
  rm -f lb.log client.log server_*.log

  # -----------------------------
  # 1) Start the Load Balancer
  # -----------------------------
  echo "Starting Load Balancer on port $LB_PORT with policy=$POLICY..."
  go run server/lb.go --port=$LB_PORT --policy=$POLICY &
  LB_PID=$!
  sleep 2

  # -----------------------------
  # 2) Start Multiple Servers
  # -----------------------------
  echo "Spawning $NUM_SERVERS servers..."
  for i in $(seq 1 $NUM_SERVERS); do
    PORT=$((START_PORT + i - 1))
    SERVER_ID="server$i"
    echo "  Starting $SERVER_ID on port $PORT..."
    go run server/server.go --port=$PORT --id=$SERVER_ID &
    sleep 1
  done

  sleep 5  # Allow servers to register and start reporting load

  # -----------------------------
  # 3) Launch Multiple Clients
  # -----------------------------
  echo "Launching $NUM_CLIENTS clients..."
  for i in $(seq 1 $NUM_CLIENTS); do
    echo "  Starting client $i with $CLIENT_REQUESTS requests, concurrency=$CLIENT_CONCURRENCY..."
    go run client/client.go \
      --lb_addr="localhost:$LB_PORT" \
      --policy="$POLICY" \
      --requests=$CLIENT_REQUESTS \
      --concurrency=$CLIENT_CONCURRENCY &
    sleep 0.5
  done

  echo "Test running for $TEST_DURATION seconds for policy $POLICY..."
  sleep $TEST_DURATION

  # -----------------------------
  # 4) Shutdown All Processes
  # -----------------------------
  echo "Terminating processes for policy $POLICY..."
  pkill -f "go run server/lb.go"
  pkill -f "go run server/server.go"
  pkill -f "go run client/client.go"
  sleep 5  # Allow processes to terminate and logs to flush

  # -----------------------------
  # 5) Analyze Logs & Record Metrics
  # -----------------------------
  # Count total responses in client.log (assuming each successful response contains the string)
  RESPONSES=$(grep -c "Response from server:" client.log)
  THROUGHPUT=$(echo "scale=2; $RESPONSES / $TEST_DURATION" | bc)

  echo "----------------------------------------" | tee -a "$FINAL_LOG"
  echo "Performance Metrics for Policy: $POLICY" | tee -a "$FINAL_LOG"
  echo "Total responses: $RESPONSES" | tee -a "$FINAL_LOG"
  echo "Throughput (responses/sec): $THROUGHPUT" | tee -a "$FINAL_LOG"
  echo "" | tee -a "$FINAL_LOG"
  echo "Load Balancer Assignment Distribution:" | tee -a "$FINAL_LOG"
  # Use grep with regex to capture both "Assigned server:" lines and policy-specific lines:
  grep -E "(Assigned server:|Round Robin -> Selected server:|Least Load -> Selected server:)" lb.log | awk '{print $NF}' | sort | uniq -c | sort -nr | tee -a "$FINAL_LOG"
  echo "----------------------------------------" | tee -a "$FINAL_LOG"
  echo "" | tee -a "$FINAL_LOG"

  sleep 5  # Pause before next test
done

echo "Scale testing complete. Final performance metrics:"
cat "$FINAL_LOG"
