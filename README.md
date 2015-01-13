## About embd-go tool

[![License: MIT.](https://img.shields.io/badge/license-MIT-orange.svg)](http://choosealicense.com/licenses/mit/)

embd-go is an embeddable command-line tool for embedding data files in Go source code, ***specially crafted for easy use with [`go generate`](http://blog.golang.org/generate)***.

## Easy use with [`go generate`](http://blog.golang.org/generate)

The embd.go is a single, self-contained, MIT-licensed, go-runnable file, so you can ***copy it verbatim*** into your own project's repository, and keep it there forever:

    Windows: copy %GOPATH%\github.com\akavel\embd-go\embd.go tools\embd.go
    Linux:   cp $GOPATH/github.com/akavel/embd-go/embd.go tools/embd.go
    
and then call from e.g. [`go generate`](http://blog.golang.org/generate) via `go run`, by putting a line like shown below in one of your Go source files:

    //go:generate go run tools/embd.go -o embd/data.go -p embd MY_DATA/HELLO.DAT


## Usage

```
USAGE: go run embd.go [FLAGS] PATH...
  -o="embd/data.go": Path to generated file.
  -p="embd": Package that the generated file should be in.
```

## Example

```
C:> go get github.com/akavel/embd-go
C:> echo Hello> hello.txt
C:> embd-go.exe -o hello.go -p main hello.txt
C:> type hello.go
// DO NOT EDIT BY HAND
//
// Generated with:
//
//      embd-go "-o" "hello.go" "-p" "main" "hello.txt"

package main

// File_hello_txt contains contents of "hello.txt" file.
var File_hello_txt = []byte{"" +
        "Hello\r\n" +
        ""}
```

