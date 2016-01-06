/*
The MIT License (MIT)

Copyright (c) 2015 Mateusz Czapliński <czapkofan@gmail.com>
Copyright (c) 2016 alxmsl <alexey.y.maslov@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

const version = "embd-go version 1.1 (2016-01-04) https://github.com/akavel/embd-go"

var (
	out = flag.String("o", "embd/data.go", "Path to generated file.")
	pkg = flag.String("p", "embd", "Package that the generated file should be in.")
	// TODO(akavel): support gzipping & unzipping when requested via option
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: %s [FLAGS] PATH...\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, version)
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	//TODO: for each path, detect if it's a file or directory
	//TODO: bail out if any normalized variable name repeats itself
	//TODO: generate file contents, converted appropriately to escaped multiline string blobs
	//TODO[later]: build only one big string constant in init(), then make the variables contain its subslices
	//TODO[later]: make sure we don't follow symlinks (for simplicity)

	contents := Contents{
		// TODO(akavel): quote them properly for command line, not via {{printf "%q"}}
		Args:  os.Args[1:],
		Pkg:   *pkg,
		Files: map[string]File{},
		Dirs:  map[string]map[string]File{},
	}
	for _, path := range flag.Args() {
		path := filepath.ToSlash(path)
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			dir, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}

			k := "Dir" + normalize(path)
			if _, exists := contents.Dirs[k]; exists {
				return fmt.Errorf("directory %s was resolved to variable %s. But last's already found", path, info.Name())
			}

			files := map[string]File{}
			for _, info := range dir {
				f, err := NewFile(path + "/" + info.Name())
				if err != nil {
					return err
				}
				files[f.VarName] = f
			}
			contents.Dirs[k] = files
		} else {
			f, err := NewFile(path)
			if err != nil {
				return err
			}
			if old, exists := contents.Files[f.VarName]; exists {
				return fmt.Errorf(
					"generated variable name conflict: '%s' resolves to"+
						"the same variable name %s as '%s'",
					f.Path, f.VarName, old.Path)
			}

			k := "File" + normalize(path)
			contents.Files[k] = f
		}
	}

	err := os.MkdirAll(filepath.Dir(*out), 0777)
	if err != nil {
		return err
	}

	w := bytes.Buffer{}
	err = template.Must(template.New("Contents").Parse(Template)).Execute(&w, contents)
	if err != nil {
		return err
	}

	code, err := format.Source(w.Bytes())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(*out, code, 0644)
	if err != nil {
		return err
	}

	return nil
}

func normalize(path string) string {
	return Normalize.ReplaceAllString("_"+filepath.Base(path), "_")
}

var Normalize = regexp.MustCompile(`[^A-Za-z0-9]+`)

func GoEscaped(buf []byte) string {
	s := fmt.Sprintf("%q", string(buf))
	return s[1 : len(s)-1]
}

func NewFile(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		return File{}, err
	}

	ch := make(chan string)
	go func() {
		defer f.Close()

		r := bufio.NewReader(f)

		buf := [20]byte{}
		for {
			n, err := io.ReadFull(r, buf[:])
			switch err {
			case io.ErrUnexpectedEOF:
				ch <- GoEscaped(buf[:n])
				fallthrough
			case io.EOF:
				close(ch)
				return
			case nil:
				ch <- GoEscaped(buf[:])
			default:
				panic(fmt.Errorf("%s: %s", path, err))
			}
		}
	}()
	return File{
		Path:          path,
		VarName:       Normalize.ReplaceAllString("File_"+filepath.Base(path), "_"),
		DataFragments: ch,
	}, nil
}

type Contents struct {
	Args  []string
	Pkg   string
	Files map[string]File
	Dirs  map[string]map[string]File
}

type File struct {
	VarName, Path string
	DataFragments <-chan string
}

var Template = `
// DO NOT EDIT BY HAND
//
// Generated with:
//
//	embd-go{{range .Args}}{{. | printf " %q"}}{{end}}

package {{.Pkg}}
{{range .Files}}
// {{.VarName}} contains contents of "{{.Path}}" file.
var {{.VarName}} = []byte("{{range .DataFragments}}{{.}}{{end}}")
{{end}}

{{range $index, $element := .Dirs}}
var {{$index}} = struct {
{{range $element}}
    {{.VarName}} []byte
{{end}}
}{
{{range $element}}
    []byte("{{range .DataFragments}}{{.}}{{end}}"),
{{end}}
}
{{end}}`[1:]
