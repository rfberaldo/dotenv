# dotenv

[![Tests Status](https://github.com/rfberaldo/dotenv/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/rfberaldo/dotenv/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rfberaldo/dotenv)](https://goreportcard.com/report/github.com/rfberaldo/dotenv)
[![Go Reference](https://pkg.go.dev/badge/github.com/rfberaldo/dotenv.svg)](https://pkg.go.dev/github.com/rfberaldo/dotenv)

> The most stupid simple dotenv package.

Multiline is not allowed, dotenv considers the key being what's before the first equal sign, and the value being what's after up to the line break.
No quotes or escapes needed.

## Install

```bash
go get github.com/rfberaldo/dotenv
```

## Usage

Load accepts multiple filepaths, if none is given, ".env" is used.

```go
if err := dotenv.Load(); err != nil {
  log.Fatal(err)
}
```

Variables are expanded by default, meaning that if your file has:

```ini
PORT=8000
HOST=localhost:${PORT}
URL=http://${HOST}
```

After loading, `os.Getenv("URL")` would produce `http://localhost:8000`.

### Overriding

By default, dotenv won't overwrite existing environment variables, the first file to set a variable will have priority.

To make it override, use `SetOverride`. In this mode, the last file to set a variable will have priority.

```go
dotenv.SetOverride(true)
if err := dotenv.Load(".env", ".env.production"); err != nil {
  log.Fatal(err)
}
```

By default, dotenv won't return error if file do not exists.
To force a file to exist, use `SetRequireFileExists`.

```go
dotenv.SetRequireFileExists(true)
if err := dotenv.Load(); err != nil {
  log.Fatal(err)
}
```

### Reading

It's also possible to get a key-value pair map instead of populating the environment.

Read follows the same rules as Load, but will return a map.
Predefined vars are only considered for expanding.

```go
kv, err := dotenv.Read() // returns map[string]string
if err != nil {
  log.Fatal(err)
}
```

### Testing

dotenv exposes a helper function for testing called `LoadTesting`.
It will use `t.Setenv()` instead of `os.Setenv`.

```go
if err := dotenv.LoadTesting(t); err != nil {
  t.Fatal(err)
}
```

### Parsing

Parse is a primitive, it receives a string and returns a key-value pair map.
However will not expand the variables.

```go
str := "PORT=8000\nHOST=localhost:${PORT}"
kv, err := dotenv.Parse(strings.NewReader(str))
// Output: map[string]string{
//   "PORT": "8000",
//   "HOST": "localhost:${PORT}",
// }
```
