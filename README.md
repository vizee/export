# export

A lite package for export variable.

## Install

```text
go install github.com/vizee/export
```

## Usage

```go
package main

import (
    "net"

    "github.com/vizee/export"
)

type StrVar string

func (v *StrVar) Get() []byte {
    return []byte(*v)
}

func (v *StrVar) Set(s string) []byte {
    o := []byte(*v)
    *v = StrVar(s)
    return o
}

func main() {
    var (
        b = true
        i = 126
        s = "str"
        a StrVar
    )
    export.Bool("bool", &b)
    export.Int("int", &i)
    export.String("str", &s)
    export.Any("any", &a)
    l, err := net.Listen("tcp", ":0")
    if err != nil {
        panic(err)
    }
    println(l.Addr().String())
    // > telnet address port
    for {
        conn, err := l.Accept()
        if err != nil {
            panic(err)
        }
        err = export.Serve(conn, conn)
        conn.Close()
        if err != nil {
            println(err)
        }
    }
}
```

## Commands

* KEYS
```text
list all keys
```

* GET <key>
```text
get value by key
```

* SET <key> <value>
```text
update value by key and return previous value
```
