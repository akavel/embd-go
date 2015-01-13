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
//      go run "-o" "hello.go" "-p" "main" "hello.txt"

package main

// File_hello_txt contains contents of "hello.txt" file.
var File_hello_txt = []byte{"" +
        "Hello\r\n" +
        ""}
```

## Single file

The embd.go is a single, self-contained, go-runnable file, so you can *copy it verbatim* into your own project's repository:

    Windows: copy %GOPATH%\github.com\akavel\embd-go\embd.go tools\embd.go
    Linux:   cp $GOPATH/github.com/akavel/embd-go/embd.go tools/embd.go
    
and then call from e.g. `go generate` via `go run`, as in:

    //go:generate go run tools/embd.go -o embd/data.go -p embd MY_DATA.dat
