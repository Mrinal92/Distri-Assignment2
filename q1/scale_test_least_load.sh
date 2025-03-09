#!/bin/bash
#
# scale_test_least_load.sh
#
# This script tests the legitimacy of the "least_load" policy.
#
# It does the following:
#   1. Starts the Load Balancer (LB) with the least_load policy.
#   2. Spawns 5 servers (on ports 9100-9104).
#   3. Launches 10 clients:
#         - Odd-indexed clients send heavy tasks (--task_type=heavy)
#         - Even-indexed clients send regular tasks (--task_type=regular)
#   4. Monitors the LB log until all servers report 0.00 load for 5 seconds,
#      then terminates the test.
#   5. Computes performance metrics (total responses, throughput, LB assignment distribution)
#      and writes a final summary log.
#
# Parameters:
NUM_SERVERS=5                # Number of servers to spawn
START_PORT=9100              # Starting port for servers (9100, 9101, …)
NUM_CLIENTS=10               # Number of client processes to launch
CLIENT_REQUESTS=20           # Each client sends this many requests
CLIENT_CONCURRENCY=2         # Number of concurrent goroutines per client
MAX_DURATION=300             # Maximum test duration in seconds (e.g., 5 minutes)
FINAL_LOG="final_least_load_performance.log"

# Clean up previous final log and logs.
rm -f "$FINAL_LOG"
rm -f lb.log client.log server_*.log

# -----------------------------
# 1) Start the Load Balancer
# -----------------------------
echo "Starting Load Balancer on port 9000 with policy=least_load..."
go run server/lb.go --port=9000 --policy=least_load &
LB_PID=$!
sleep 2

# -----------------------------
# 2) Spawn Multiple Servers
# -----------------------------
echo "Spawning $NUM_SERVERS servers..."
for i in $(seq 1 $NUM_SERVERS); do
    PORT=$((START_PORT + i - 1))
    SERVER_ID="server$i"
    echo "  Starting $SERVER_ID on port $PORT..."
    go run server/server.go --port=$PORT --id=$SERVER_ID &
    sleep 1
done
sleep 5  # Allow servers to register and start reporting load.

# -----------------------------
# 3) Launch Multiple Clients (Mixed Task Types)
# -----------------------------
echo "Launching $NUM_CLIENTS clients..."
for i in $(seq 1 $NUM_CLIENTS); do
    if (( i % 2 == 1 )); then
        TASK_TYPE="heavy"
    else
        TASK_TYPE="regular"
    fi
    echo "  Starting client $i: $CLIENT_REQUESTS requests, concurrency=$CLIENT_CONCURRENCY, task_type=$TASK_TYPE..."
    go run client/client.go \
      --lb_addr="localhost:9000" \
      --policy="least_load" \
      --requests=$CLIENT_REQUESTS \
      --concurrency=$CLIENT_CONCURRENCY \
      --task_type="$TASK_TYPE" &
    sleep 0.5
done

# -----------------------------
# 4) Monitor LB Load until Idle for 5 Seconds or MAX_DURATION reached
# -----------------------------
echo "Monitoring LB load for idle state (all servers reporting 0.00 load for 5 sec)..."
IDLE_FOUND=false
ELAPSED=0
while [ $ELAPSED -lt $MAX_DURATION ]; do
    sleep 5
    ELAPSED=$(( ELAPSED + 5 ))
    ALL_IDLE=true
    for i in $(seq 1 $NUM_SERVERS); do
        LATEST=$(grep "Updated load for server server$i:" lb.log | tail -n 1)
        if [[ $LATEST != *"0.00"* ]]; then
            ALL_IDLE=false
            break
        fi
    done
    if $ALL_IDLE; then
        echo "All servers idle for the last 5 seconds."
        IDLE_FOUND=true
        break
    else
        echo "Not all servers idle yet (elapsed: $ELAPSED sec)..."
    fi
done

# -----------------------------
# 5) Shutdown All Processes
# -----------------------------
echo "Terminating all LB, server, and client processes..."
pkill -f "go run server/lb.go"
pkill -f "go run server/server.go"
pkill -f "go run client/client.go"
sleep 5

# -----------------------------
# 6) Analyze Logs & Generate Final Performance Report
# -----------------------------
RESPONSES=$(grep -c "Response from server:" client.log)
if [ $ELAPSED -gt 0 ]; then
    THROUGHPUT=$(echo "scale=2; $RESPONSES / $ELAPSED" | bc)
else
    THROUGHPUT=0
fi

{
    echo "----------------------------------------"
    echo "Final Performance Metrics for LB Policy: least_load"
    echo "Test duration (seconds): $ELAPSED"
    echo "Total responses: $RESPONSES"
    echo "Throughput (responses/sec): $THROUGHPUT"
    echo ""
    echo "Load Balancer Assignment Distribution:"
    # Capture lines containing "Assigned server:" and policy-specific messages.
    grep -E "(Assigned server:|Round Robin -> Selected server:|Least Load -> Selected server:)" lb.log | awk '{print $NF}' | sort | uniq -c | sort -nr
    echo "----------------------------------------"
} >> "$FINAL_LOG"

echo "Scale test complete. Final performance report:"
cat "$FINAL_LOG"
