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
C:> go run $GOPATH/github.com/akavel/embd-go/embd.go -o hello.go -p main hello.txt
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
