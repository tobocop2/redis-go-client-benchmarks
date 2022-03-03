redis:
	@docker pull redis >/dev/null 2>&1 && docker run --name dev-redis -p 6379:6379  -d redis >/dev/null 2>&1 || true

build-go-redis-benchmark: go-redis
	@cd $^ && go build

build-redigo-benchmark: redigo
	@cd $^ && go build

bench-go-redis: redis build-go-redis-benchmark
	ulimit -n 99999; ./go-redis/go-redis

bench-redigo: redis build-redigo-benchmark
	ulimit -n 99999; ./redigo/redigo

clean:
	rm -f ./go-redis/go-redis ./redigo/redigo
