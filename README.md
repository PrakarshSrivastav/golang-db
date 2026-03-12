# go-filedb

A lightweight, file-based JSON database written in Go. No external database server required - records are stored as `.json` files on disk, organised into collections, with full concurrency safety via per-collection mutexes.

---

## Features

- **Zero dependencies** (beyond a logger) — no database server, no setup
- **Concurrent-safe** reads and writes using per-collection mutex locking
- **Atomic writes** via temp file + rename pattern to prevent data corruption
- **Simple CRUD API** — Write, Read, ReadAll, Delete
- **Human-readable storage** — all records stored as indented `.json` files

---

## Installation

```bash
go get github.com/jcelliott/lumber
```

Clone the repo and run:

```bash
go run main.go
```

---

## Usage

### Initialise the database

```go
db, err := New("./mydata", nil)
if err != nil {
    log.Fatal(err)
}
```

The second argument accepts an `*Options` struct for a custom logger. Pass `nil` to use the default console logger.

---

### Write a record

```go
type User struct {
    Name    string
    Age     json.Number
    Contact string
    Company string
    Address Address
}

db.Write("users", "John", User{
    Name:    "John",
    Age:     "23",
    Contact: "22333445",
    Company: "Golden State",
    Address: Address{"Lucknow", "Uttar Pradesh", "India", "441003"},
})
```

Records are saved to `<dir>/<collection>/<resource>.json`.

---

### Read a single record

```go
var user User
err := db.Read("users", "John", &user)
```

---

### Read all records in a collection

```go
records, err := db.ReadAll("users")
for _, r := range records {
    var user User
    json.Unmarshal([]byte(r), &user)
    fmt.Println(user)
}
```

---

### Delete a record

```go
err := db.Delete("users", "John") // deletes users/John.json
```

To delete an entire collection:

```go
err := db.Delete("users", "") // removes the users/ directory
```

---

## File Structure

```
./
└── users/
    ├── John.json
    ├── Paul.json
    ├── Curry.json
    └── ...
```

Each file contains a single indented JSON object:

```json
{
    "Name": "John",
    "Age": "23",
    "Contact": "22333445",
    "Company": "Golden State",
    "Address": {
        "City": "Lucknow",
        "State": "Uttar Pradesh",
        "Country": "India",
        "Pincode": "441003"
    }
}
```

---

## Concurrency Model

- A **global mutex** protects the `mutexes` map itself
- A **per-collection mutex** is created lazily on first access
- Writes use an **atomic rename** (`file.tmp` → `file.json`) to prevent partial writes on crash

---

## Limitations

- Not suitable for high-throughput production workloads — designed for lightweight local storage
- No query/filter support — reads return raw JSON strings
- No indexing — `ReadAll` scans all files in the collection directory
- `ioutil` functions used are deprecated in Go 1.16+; consider migrating to `os` and `io` packages

---

## License

MIT
