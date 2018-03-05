Yet Another Data Embedder
=========================

This tool will generate a go file with an embedded gzipped tarball string.
On your program startup, it will uncompress itself in memory.

## Installation

```
go get -u github.com/tuxlinuxien/yade
```

## Usage

```
$> yade -h
Usage of yade:
  -dest string
        destination file (default "emb.go")
  -package string
        package name (default "emb")
  -src string
        source files or directory (default "./")
```


Example

```
$> yade -src ui/static/,ui/index.html,ui/dest
```

## Attribution

This package is highly inspired from [go-bindata](https://github.com/jteeuwen/go-bindata)
