for i in $(seq 1 10000); do
  curl -X POST http://localhost:8080/post -d 'hello'
done