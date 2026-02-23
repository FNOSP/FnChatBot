# Table-Driven Tests

The idiomatic Go pattern for writing comprehensive, maintainable tests using data-driven test cases.

## Basic Pattern

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed signs", -2, 3, 1},
        {"zeros", 0, 0, 0},
        {"large numbers", 1000000, 2000000, 3000000},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

## With Error Testing

```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name      string
        a, b      float64
        expected  float64
        expectErr bool
    }{
        {"normal division", 10, 2, 5, false},
        {"negative result", -10, 2, -5, false},
        {"division by zero", 10, 0, 0, true},
        {"zero numerator", 0, 5, 0, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)

            if tt.expectErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if result != tt.expected {
                t.Errorf("Divide(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

## Parallel Execution

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"lowercase", "hello", "HELLO"},
        {"uppercase", "WORLD", "WORLD"},
        {"mixed", "HeLLo", "HELLO"},
        {"empty", "", ""},
        {"numbers", "123", "123"},
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            got := strings.ToUpper(tt.input)
            if got != tt.want {
                t.Errorf("ToUpper(%q) = %q, want %q", tt.input, got, tt.want)
            }
        })
    }
}
```

## Complex Input Types

```go
func TestProcessUser(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        options ProcessOptions
        want    ProcessResult
        wantErr bool
    }{
        {
            name: "standard user",
            user: User{Name: "John", Age: 30},
            options: ProcessOptions{Validate: true},
            want: ProcessResult{Status: "valid"},
            wantErr: false,
        },
        {
            name: "invalid age",
            user: User{Name: "Jane", Age: -1},
            options: ProcessOptions{Validate: true},
            want: ProcessResult{},
            wantErr: true,
        },
        {
            name: "skip validation",
            user: User{Name: "", Age: 0},
            options: ProcessOptions{Validate: false},
            want: ProcessResult{Status: "skipped"},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessUser(tt.user, tt.options)

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ProcessUser() = %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

## With Setup/Teardown

```go
func TestDatabase(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*sql.DB)
        query   string
        want    []User
        wantErr bool
    }{
        {
            name: "find all users",
            setup: func(db *sql.DB) {
                db.Exec("INSERT INTO users (name) VALUES ('Alice'), ('Bob')")
            },
            query: "SELECT * FROM users",
            want: []User{{Name: "Alice"}, {Name: "Bob"}},
            wantErr: false,
        },
        {
            name: "empty table",
            setup: func(db *sql.DB) {},
            query: "SELECT * FROM users",
            want: []User{},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db := setupTestDB(t)
            tt.setup(db)

            got, err := queryUsers(db, tt.query)

            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }

            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

## Edge Cases Checklist

Always include these edge cases:

```go
func TestComprehensive(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {"empty", "", ""},
        {"single char", "a", "A"},
        {"whitespace", "   ", "   "},
        {"unicode", "héllo wörld", "HÉLLO WÖRLD"},
        {"emoji", "hello 👋", "HELLO 👋"},
        {"very long", strings.Repeat("a", 10000), strings.Repeat("A", 10000)},
        {"newlines", "hello\nworld", "HELLO\nWORLD"},
        {"tabs", "hello\tworld", "HELLO\tWORLD"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := strings.ToUpper(tt.input)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

## Test Helpers

```go
func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func assertError(t *testing.T, err error, wantErr bool) {
    t.Helper()
    if wantErr && err == nil {
        t.Error("expected error, got nil")
    }
    if !wantErr && err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

func TestWithHelpers(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"add", 2, 3, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            assertEqual(t, got, tt.want)
        })
    }
}
```

## Best Practices

1. **Always use `t.Run()`** - Creates subtests with clear names
2. **Capture range variable** - Use `tt := tt` for parallel tests
3. **Name test cases descriptively** - Documents what's being tested
4. **Include edge cases** - Empty, nil, max, min, invalid
5. **Test errors explicitly** - Use `wantErr` boolean
6. **Use `t.Helper()`** - Improves error line numbers
7. **Keep tests independent** - No shared state between tests
