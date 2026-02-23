# Fuzzing

Go fuzzing (Go 1.18+) for discovering edge cases and security vulnerabilities through random input generation.

## Basic Fuzzing

```go
func FuzzReverse(f *testing.F) {
    testcases := []string{"hello", "world", "123", ""}
    for _, tc := range testcases {
        f.Add(tc)
    }

    f.Fuzz(func(t *testing.T, input string) {
        reversed := Reverse(input)
        doubleReversed := Reverse(reversed)

        if input != doubleReversed {
            t.Errorf("Reverse(Reverse(%q)) = %q, want %q", input, doubleReversed, input)
        }
    })
}
```

Run with:
```bash
go test -fuzz=FuzzReverse
go test -fuzz=FuzzReverse -fuzztime=30s
go test -fuzz=. -fuzztime=1m
```

## Multiple Parameters

```go
func FuzzAdd(f *testing.F) {
    f.Add(1, 2)
    f.Add(0, 0)
    f.Add(-1, 1)
    f.Add(int64(9223372036854775807), int64(1))

    f.Fuzz(func(t *testing.T, a, b int64) {
        result := Add(a, b)

        if a > 0 && b > 0 && result < 0 {
            t.Errorf("Add(%d, %d) = %d; potential overflow", a, b, result)
        }
    })
}
```

## Supported Types

```go
func FuzzTypes(f *testing.F) {
    f.Add("string", 42, true, 3.14, []byte("bytes"))

    f.Fuzz(func(t *testing.T, s string, i int, b bool, f float64, data []byte) {
        _ = s
        _ = i
        _ = b
        _ = f
        _ = data
    })
}
```

Supported types: `string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `[]byte`.

## Property-Based Testing

```go
func FuzzJSON(f *testing.F) {
    f.Add(`{"name": "test", "value": 42}`)

    f.Fuzz(func(t *testing.T, jsonStr string) {
        var data map[string]interface{}

        if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
            return
        }

        reencoded, err := json.Marshal(data)
        if err != nil {
            t.Errorf("failed to re-encode: %v", err)
            return
        }

        var data2 map[string]interface{}
        if err := json.Unmarshal(reencoded, &data2); err != nil {
            t.Errorf("failed to unmarshal re-encoded: %v", err)
        }
    })
}
```

## Finding Edge Cases

```go
func FuzzParseInt(f *testing.F) {
    f.Add("123")
    f.Add("-456")
    f.Add("0")

    f.Fuzz(func(t *testing.T, s string) {
        val, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return
        }

        str := strconv.FormatInt(val, 10)
        val2, err := strconv.ParseInt(str, 10, 64)
        if err != nil {
            t.Errorf("round-trip failed: %v", err)
        }

        if val != val2 {
            t.Errorf("round-trip mismatch: %d != %d", val, val2)
        }
    })
}
```

## Fuzzing with Custom Types

```go
type User struct {
    Name  string
    Email string
    Age   int
}

func FuzzUserValidation(f *testing.F) {
    f.Add("John", "john@example.com", 25)

    f.Fuzz(func(t *testing.T, name, email string, age int) {
        user := User{Name: name, Email: email, Age: age}
        err := ValidateUser(user)

        if age < 0 && err == nil {
            t.Error("expected error for negative age")
        }
        if age >= 0 && age <= 150 && err != nil {
            t.Errorf("unexpected error for valid age: %v", err)
        }
    })
}
```

## Crashers and Findings

When fuzzing finds a crasher, it saves it:

```
testdata/fuzz/FuzzReverse/abc123
```

Contents:
```
go test fuzz v1
string("crasher input")
```

To reproduce:
```bash
go test -run=FuzzReverse/abc123
```

## Fuzzing Options

```bash
go test -fuzz=FuzzName                    # Run specific fuzz
go test -fuzz=.                           # Run all fuzz tests
go test -fuzztime=30s                     # Fuzz for 30 seconds
go test -fuzztime=1000x                   # Fuzz for 1000 iterations
go test -fuzzminimizedtimeout=5s          # Timeout for minimizing
go test -fuzz=FuzzName -run=^$            # Only fuzz, skip unit tests
```

## Best Practices

1. Add seed corpus with known inputs
2. Test properties, not exact outputs
3. Handle invalid inputs gracefully
4. Use fuzzing alongside unit tests
5. Commit crashers to prevent regressions
6. Keep fuzz functions simple and fast

## Example: URL Parser

```go
func FuzzURLParse(f *testing.F) {
    f.Add("https://example.com/path?query=value")
    f.Add("http://localhost:8080")
    f.Add("/relative/path")

    f.Fuzz(func(t *testing.T, urlStr string) {
        u, err := url.Parse(urlStr)
        if err != nil {
            return
        }

        reencoded := u.String()
        u2, err := url.Parse(reencoded)
        if err != nil {
            t.Errorf("round-trip failed: %v", err)
            return
        }

        if u.Scheme != u2.Scheme || u.Host != u2.Host {
            t.Errorf("round-trip mismatch: %+v != %+v", u, u2)
        }
    })
}
```
