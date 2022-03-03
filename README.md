README
======

This repository benchmarks go-redis and redigo, two very popular redis go clients and demonstrates
how each handle some simple operations behind a redis connection pool

# Prerequisites

* docker
* go

# Running

To see the output from the redigo benchmark

Run

```bash
make bench-redigo
```

To see the output from the go-redis benchmark

Run

```bash
make bench-go-redis
```
