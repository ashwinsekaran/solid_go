for i in $(seq 1 100); do
  curl -s -o /dev/null -w "%{http_code}\n" \
    -H "X-User-ID: ashwin" \
    http://localhost:8080/
done