#!/bin/sh
set -e

go mod download

air -c .air.recorder.toml &
AIR_PID=$!

LOG_FILE="/app/logs/recorder.log"
mkdir -p "$(dirname "$LOG_FILE")"

while [ ! -f "$LOG_FILE" ]; do
  sleep 0.5
done
tail -n 0 -F "$LOG_FILE" &
TAIL_PID=$!

cleanup() {
  kill "$TAIL_PID" 2>/dev/null || true
  kill "$AIR_PID" 2>/dev/null || true
}
trap cleanup INT TERM EXIT

wait "$AIR_PID"
