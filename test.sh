#!/bin/bash

echo "🚀 Launching 50 concurrent requests..."

# Loop 50 times
for i in {1..50}
do
   # Run curl in the background using '&'
   curl -s -X POST http://localhost:4000/v1/send \
     -H "Content-Type: application/json" \
     -d "{\"from\":\"admin@test.com\", \"to\":\"user$i@test.com\", \"subject\":\"Concurrent Go Request $i\", \"body\":\"Load test payload\"}" &
done

# Wait for all background jobs to finish
wait

echo "✅ Done! All requests sent."