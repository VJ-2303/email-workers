#!/bin/bash

echo "🚀 Launching 50 concurrent requests..."

for i in {1..15}
do
   curl -s -X POST http://localhost:4000/v1/send \
     -H "Content-Type: application/json" \
     -d "{\"from\":\"admin@test.com\", \"to\":\"vvanam37@gmail.com\", \"subject\":\"Concurrent Go Request $i\", \"body\":\"Load test payload\"}" &
done

wait

echo "Done! All requests sent."