# Dotenv

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

Variables are expanded by default, meaning that if your .env has:

```conf
HOST=localhost:8000
URL=http://${HOST}
```

After loading, `os.Getenv("URL")` will produce `http://localhost:8000`.

By default, dotenv won't overwrite existing environment variables, the first file to set a variable will have priority.

To make it override, use `SetOverride`. In this mode, the last file to set a variable will have priority.

```go
dotenv.SetOverride()
if err := dotenv.Load(".env", ".env.production"); err != nil {
  log.Fatal(err)
}
```

By default, dotenv won't return error if file do not exists.
To force a file to exist, use `SetRequireFileExists`.

```go
dotenv.SetRequireFileExists()
if err := dotenv.Load(); err != nil {
  log.Fatal(err)
}
```

## Parsing

Parse a string to get a key-value pair map. Parse will not expand the variable.

```go
kv, err := dotenv.Parse("ENV=dev\nHOST=localhost:8000")
// Output: map[string]string{
//   "ENV":  "dev",
//   "HOST": "localhost:8000",
// }
```
