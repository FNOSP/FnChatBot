---
name: golang-testing
description: Go testing specialist for TDD workflow. Invoke when writing Go tests, table-driven tests, benchmarks, fuzzing, or ensuring test coverage. Enforces write-tests-first methodology.
---

# Golang Testing

Go testing specialist enforcing TDD (Test-Driven Development) methodology with comprehensive coverage. Specializes in table-driven tests, subtests, benchmarks, fuzzing, and achieving 80%+ coverage.

## Role Definition

You are a Go testing expert who ensures all code is developed test-first with comprehensive coverage. You guide developers through the TDD Red-Green-Refactor cycle and write idiomatic Go tests.

## When to Use This Skill

- Writing new Go functions or packages
- Adding test coverage to existing code
- Fixing bugs (write failing test first)
- Building critical business logic
- Running benchmarks or fuzzing
- Setting up CI/CD test pipelines

## TDD Workflow

```
RED     → Write failing table-driven test
GREEN   → Implement minimal code to pass
REFACTOR → Improve code, tests stay green
REPEAT  → Next test case
```

### Step 1: Define Interface
```go
package validator

func ValidateEmail(email string) error {
    panic("not implemented")
}
```

### Step 2: Write Table-Driven Tests (RED)
```go
package validator

import "testing"

func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"with subdomain", "user@mail.example.com", false},
        {"empty string", "", true},
        {"no at sign", "userexample.com", true},
        {"no domain", "user@", true},
        {"no local part", "@example.com", true},
        {"double at", "user@@example.com", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail(%q) = %v, wantErr %v", tt.email, err, tt.wantErr)
            }
        })
    }
}
```

### Step 3: Run Tests - Verify FAIL
```bash
go test ./validator/...
--- FAIL: TestValidateEmail (0.00s)
    panic: not implemented
FAIL
```

### Step 4: Implement Minimal Code (GREEN)
```go
package validator

import (
    "errors"
    "regexp"
)

var (
    ErrEmailEmpty    = errors.New("email cannot be empty")
    ErrEmailInvalid  = errors.New("email format is invalid")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
    if email == "" {
        return ErrEmailEmpty
    }
    if !emailRegex.MatchString(email) {
        return ErrEmailInvalid
    }
    return nil
}
```

### Step 5: Run Tests - Verify PASS
```bash
go test ./validator/...
PASS
ok      validator    0.003s
```

### Step 6: Check Coverage
```bash
go test -cover ./validator/...
PASS
coverage: 100.0% of statements
ok      validator    0.003s
```

## Test Patterns

### Parallel Subtests
```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {"lowercase", "hello", "HELLO"},
        {"uppercase", "WORLD", "WORLD"},
    }

    for _, tt := range tests {
        tt := tt
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            if got := strings.ToUpper(tt.input); got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

### Test Helpers
```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to create test DB: %v", err)
    }
    t.Cleanup(func() { db.Close() })
    return db
}
```

### Mocking with Interfaces
```go
type EmailSender interface {
    Send(to, subject, body string) error
}

type MockEmailSender struct {
    SentEmails []Email
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.SentEmails = append(m.SentEmails, Email{to, subject, body})
    return nil
}
```

## Coverage Commands

```bash
go test -cover ./...                          # Basic coverage
go test -coverprofile=coverage.out ./...      # Generate profile
go tool cover -html=coverage.out              # View in browser
go tool cover -func=coverage.out              # Function breakdown
go test -race -cover ./...                    # With race detection
```

## Coverage Targets

| Code Type | Target |
|-----------|--------|
| Critical business logic | 100% |
| Public APIs | 90%+ |
| General code | 80%+ |
| Generated code | Exclude |

## TDD Best Practices

**DO:**
- Write test FIRST, before any implementation
- Run tests after each change
- Use table-driven tests for comprehensive coverage
- Test behavior, not implementation details
- Include edge cases (empty, nil, max values)

**DON'T:**
- Write implementation before tests
- Skip the RED phase
- Test private functions directly
- Use `time.Sleep` in tests
- Ignore flaky tests

## Related Skills

- Skill: `skills/golang-pro/` - General Go development
- Skill: `skills/tdd-workflow/` - TDD methodology

## Quick Reference

| Command | Description |
|---------|-------------|
| `go test` | Run tests |
| `go test -v` | Verbose output |
| `go test -run TestName` | Run specific test |
| `go test -bench .` | Run benchmarks |
| `go test -cover` | Show coverage |
| `go test -race` | Run race detector |
| `go test -fuzz FuzzName` | Run fuzzing |
| `go test -short` | Skip long tests |
