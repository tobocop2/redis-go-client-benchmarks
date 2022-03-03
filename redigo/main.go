package main

import (
    "log"
    "time"
    "errors"
    "github.com/gomodule/redigo/redis"
    "sync"
)

var pool *redis.Pool

// can change this to control how many redis operations to run per worker
const N = 10

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

func NewPool(poolSize int) *redis.Pool {
    opts := []redis.DialOption{
        redis.DialConnectTimeout(time.Second * 10),
    }
    pool := &redis.Pool{
        MaxActive:   poolSize,
        MaxIdle:     poolSize,
        IdleTimeout: time.Second * 60,
    }

    pool.Dial = func() (redis.Conn, error) {
        log.Println(pool.ActiveCount(), "opening new connection")

        var err error
        var conn redis.Conn
        const addr = "127.0.0.1:6379"
        for i := 0; i < 10; i++ {
            time.Sleep(time.Second * time.Duration(i))
            conn, err = redis.Dial("tcp", addr, opts...)
            if err == nil {
                return conn, nil
            }
            log.Println("cannot connect to redis, try again", i, addr, err)
        }
        log.Fatal("failed to connect to redis after 10 tries", addr, err)
        return nil, errors.New("failed to connect to redis")
    }
    return pool
}


func set(key string, val string) error {
    defer timeTrack(time.Now(), "Redis SET")
    conn := pool.Get()
    defer conn.Close()

    err := conn.Send("SET", key, val)
    if err != nil {
        log.Printf("ERROR: fail set key %s, val %s, error %s", key, val, err.Error())
        return err
    }

    return nil
}

func get(key string) (string, error) {
    conn := pool.Get()
    defer conn.Close()

    s, err := redis.String(conn.Do("GET", key))
    if err != nil {
        log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
        return "", err
    }

    return s, nil
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
