package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var embeddedPart1 = `// DO NOT MODIFY
// This file has been generated by github.com/tuxlinuxien/yade
package %s

import (
"archive/tar"
"bytes"
"compress/gzip"
"io"
"log"
"os"
)

var dataContent = "`

var embeddedPart2 = `"

var assetFS = map[string][]byte{}

// Files .
func Files() []string {
	lFiles := []string{}
	for fname := range assetFS {
		lFiles = append(lFiles, fname)
	}
	return lFiles
}

// Exists .
func Exists(fname string) bool {
	_, ok := assetFS[fname]
	return ok
}

// Open .
func Open(fname string) ([]byte, error) {
	if !Exists(fname) {
		return nil, os.ErrNotExist
	}
	return bytes.NewBuffer(assetFS[fname]).Bytes(), nil
}

func init() {
	buff := bytes.NewBufferString(dataContent)
	gz, err := gzip.NewReader(buff)
	if err != nil {
		log.Fatalln(err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}
		fileData := &bytes.Buffer{}
		io.Copy(fileData, tr)
		assetFS[header.Name] = fileData.Bytes()
	}
}
`

var (
	packageName = "emb"
	destFile    = "emb.go"
	srcFiles    = "./"
	fileList    = []string{}
)

func addFilesToList(p string) {
	finfo, err := os.Stat(p)
	if err != nil {
		log.Fatalln(err)
	}
	if finfo.IsDir() {
		filepath.Walk(p, func(subp string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalln(err)
			}
			if info.IsDir() {
				if p != subp {
					addFilesToList(subp)
				}
			} else {
				fileList = append(fileList, subp)
			}
			return nil
		})
	} else {
		fileList = append(fileList, p)
	}
}

func openFile(f string, tr *tar.Writer) {
	fmt.Fprintln(os.Stderr, "add file:"+f)
	fdata, err := os.Open(f)
	if err != nil {
		log.Fatalln(err)
	}
	defer fdata.Close()
	tHeader := &tar.Header{}
	finfo, _ := fdata.Stat()
	tHeader.Name = f
	tHeader.Size = finfo.Size()
	tHeader.ModTime = finfo.ModTime()
	tr.WriteHeader(tHeader)
	l, err := io.Copy(tr, fdata)
	if l != finfo.Size() {
		log.Fatalln("size mismatch")
	}
}

func printFile(b []byte) {
	fmt.Printf(embeddedPart1, packageName)
	out := ""
	for i := range b {
		if b[i] >= 'a' && b[i] <= 'z' {
			out += string(b[i])
		} else if b[i] >= 'A' && b[i] <= 'Z' {
			out += string(b[i])
		} else if b[i] >= '0' && b[i] <= '9' {
			out += string(b[i])
		} else {
			out += fmt.Sprintf("\\x%02x", int(b[i]))
		}
		if len(out) >= 2048 {
			fmt.Print(out)
			out = ""
		}
	}
	if len(out) > 0 {
		fmt.Print(out)
	}
	fmt.Print(embeddedPart2)
}

func main() {
	flag.StringVar(&packageName, "package", packageName, "package name")
	flag.StringVar(&destFile, "dest", destFile, "destination file")
	flag.StringVar(&srcFiles, "src", srcFiles, "source files or directory")
	flag.Parse()

	files := strings.Split(srcFiles, ",")
	for i := range files {
		addFilesToList(files[i])
	}
	buff := bytes.Buffer{}
	gr := gzip.NewWriter(&buff)
	tr := tar.NewWriter(gr)
	done := map[string]bool{}
	for i := range fileList {
		_, ok := done[fileList[i]]
		if ok {
			continue
		}
		done[fileList[i]] = true
		openFile(fileList[i], tr)
	}
	tr.Flush()
	tr.Close()
	gr.Flush()
	gr.Close()
	printFile(buff.Bytes())
}
