package main

import (
    "log"
    "sync"
    "time"
    "context"

    "github.com/go-redis/redis/v8"
)

var (
    pool *redis.Client
    ctx = context.Background()
)

// can change this to control how many redis operations to run per worker
const N = 10

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

func redisOptions() *redis.Options {
    return &redis.Options{
        Addr: "127.0.0.1:6379",
        DialTimeout:  10 * time.Second,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        MaxRetries: 10,
        PoolTimeout:        30 * time.Second,
        IdleTimeout:        time.Minute,
        IdleCheckFrequency: 100 * time.Millisecond,
    }
}


func NewPool(poolSize int) *redis.Client {
    opt := redisOptions()
    // opt.MinIdleConns = poolSize
    opt.MaxConnAge = 0
    opt.PoolSize = poolSize
    pool := redis.NewClient(opt)
    return pool
}

func set(key string, val string) error {
    defer timeTrack(time.Now(), "Redis SET")

    status := pool.Set(ctx, key, val, 0)
    if err := status.Err(); err != nil {
        log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
        return err
    }

    return nil
}

func benchmark(poolSize, exceedPoolSizeBy int) time.Duration {
    time.Sleep(time.Second)
    var wg sync.WaitGroup
    pool = NewPool(poolSize)

    log.Println("Connected!")
    var startTime time.Time = time.Now()

    // exceeding the amount of workers in the pool to demonstrate
    for i := 0; i < poolSize + exceedPoolSizeBy; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 0; i < N; i++ {
                set("foo", "bar")
            }
        }()
    }
    wg.Wait()
    return time.Since(startTime)
}

 func main() {
    var totalDurationA time.Duration
    var totalDurationB time.Duration
    var totalDurationC time.Duration
    var totalDurationD time.Duration
    var totalDurationE time.Duration
    var totalDurationF time.Duration

    const n = 10
    for i := 0; i < n; i++ {
         totalDurationA += benchmark(500, 0)
         totalDurationB += benchmark(250, 0)
         totalDurationC += benchmark(125, 0)
         totalDurationD += benchmark(62, 0)
         totalDurationE += benchmark(31, 0)
         totalDurationF += benchmark(15, 0)
     }
     log.Println()
     log.Printf("Total average duration for pool with size %d: %s", 500, totalDurationA / n)
     log.Printf("Total average duration for pool with size %d: %s", 250, totalDurationB / n)
     log.Printf("Total average duration for pool with size %d: %s", 125, totalDurationC / n)
     log.Printf("Total average duration for pool with size %d: %s", 62, totalDurationD / n)
     log.Printf("Total average duration for pool with size %d: %s", 31, totalDurationE / n)
     log.Printf("Total average duration for pool with size %d: %s", 15, totalDurationF / n)

     log.Println("Demonstrating exceeding max pool size by 1")
     benchmark(1, 1)
 }
