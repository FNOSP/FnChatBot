# Benchmarks

Go benchmarking with `testing.B` for performance measurement and optimization.

## Basic Benchmark

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(100, 200)
    }
}
```

Run with:
```bash
go test -bench .
go test -bench=BenchmarkAdd
go test -bench=. -benchmem
```

## Benchmark with Subtests

```go
func BenchmarkStringOperations(b *testing.B) {
    benchmarks := []struct {
        name  string
        input string
    }{
        {"short", "hello"},
        {"medium", strings.Repeat("hello", 10)},
        {"long", strings.Repeat("hello", 100)},
    }

    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                _ = strings.ToUpper(bm.input)
            }
        })
    }
}
```

## Benchmark with Setup

```go
func BenchmarkMapAccess(b *testing.B) {
    m := make(map[string]int)
    for i := 0; i < 1000; i++ {
        m[fmt.Sprintf("key%d", i)] = i
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _ = m["key500"]
    }
}

func BenchmarkWithStopTimer(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        data := generateTestData()
        b.StartTimer()

        process(data)
    }
}
```

## Parallel Benchmark

```go
func BenchmarkConcurrentAccess(b *testing.B) {
    var counter int64

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            atomic.AddInt64(&counter, 1)
        }
    })
}

func BenchmarkMutexVsChannel(b *testing.B) {
    b.Run("mutex", func(b *testing.B) {
        var mu sync.Mutex
        var counter int
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                mu.Lock()
                counter++
                mu.Unlock()
            }
        })
    })

    b.Run("channel", func(b *testing.B) {
        ch := make(chan struct{}, 1)
        var counter int
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                ch <- struct{}{}
                counter++
                <-ch
            }
        })
    })
}
```

## Memory Allocation Benchmark

```go
func BenchmarkAllocation(b *testing.B) {
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        s := make([]int, 1000)
        _ = s
    }
}

func BenchmarkSliceAppend(b *testing.B) {
    b.Run("with-capacity", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            s := make([]int, 0, 1000)
            for j := 0; j < 1000; j++ {
                s = append(s, j)
            }
        }
    })

    b.Run("without-capacity", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            s := make([]int, 0)
            for j := 0; j < 1000; j++ {
                s = append(s, j)
            }
        }
    })
}
```

## Comparing Implementations

```go
func BenchmarkJSON(b *testing.B) {
    data := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }

    b.Run("encoding/json", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            json.Marshal(data)
        }
    })

    b.Run("json-iterator", func(b *testing.B) {
        b.ReportAllocs()
        for i := 0; i < b.N; i++ {
            jsoniter.Marshal(data)
        }
    })
}
```

## CPU and Memory Profiling

```bash
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof

go test -bench=. -memprofile=mem.prof
go tool pprof mem.prof

go test -bench=. -blockprofile=block.prof
go test -bench=. -mutexprofile=mutex.prof
```

## Benchmark Output

```
BenchmarkAdd-8           1000000000          0.256 ns/op
BenchmarkSliceAppend-8      500000          3245 ns/op        8192 B/op          1 allocs/op
```

| Column | Description |
|--------|-------------|
| `ns/op` | Nanoseconds per operation |
| `B/op` | Bytes allocated per operation |
| `allocs/op` | Allocations per operation |

## Best Practices

1. Use `b.ResetTimer()` to exclude setup time
2. Use `b.ReportAllocs()` to track memory allocations
3. Run multiple times with `-count=N` for stable results
4. Compare before/after when optimizing
5. Use `b.RunParallel()` for concurrent benchmarks
