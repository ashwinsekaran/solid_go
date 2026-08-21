#!/bin/bash
# batch.sh — POST 1000 users one at a time to the existing /users endpoint

URL="http://localhost:8080/users"
COUNT=1000

for i in $(seq 1 $COUNT); do
  curl -s -X POST "$URL" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer abcxyz" \
    -d "{\"id\":\"$i\",\"name\":\"user-$i\",\"email\":\"user$i@example.com\"}" \
    > /dev/null
  echo "posted $i"
done

echo "done — posted $COUNT users"