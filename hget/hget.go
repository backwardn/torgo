package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"log"
	"io"
	"net/url"
	"path"
)

import (
	"github.com/as/mute"
)

const Prefix = "hget: "
var args struct {
	h, q bool
	v bool
	a bool
}
var f *flag.FlagSet


func main() {
	out := io.Writer(os.Stdout)
	if args.h || args.q {
		usage(); os.Exit(0)
	}
	resp, err := http.Get(f.Args()[0])
	if err != nil{
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if args.a {
		u, err := url.Parse(f.Args()[0])
		_, file := path.Split(u.Path)
		fd, err := os.Create(file)
		no(err)
		defer fd.Close()
		out = fd
	}
	if _, err = io.Copy(out, resp.Body); err != nil{
		log.Fatalln(err)
	}
}

func no(err error){
	if err != nil{
		log.Fatalln(err)
	}
}

func init() {
	f = flag.NewFlagSet("main", flag.ContinueOnError)
	f.BoolVar(&args.v, "v", false, "")
	f.BoolVar(&args.h, "h", false, "")
	f.BoolVar(&args.q, "?", false, "")
	f.BoolVar(&args.a, "a", false, "")
	err := mute.Parse(f, os.Args[1:])
	if err != nil {
		printerr(err)
		os.Exit(1)
	}
}
func usage() {
	fmt.Println(`
NAME
	hget - http get url

SYNOPSIS
	hget url

DESCRIPTION
	Hget writes the contents of an http url to stdout

	Options:
	-a    auto name file and save to current directory
	-v    verbose, print response headers to stderr

EXAMPLE
	hget https://downover.io > downover.html

`)
}

func println(v ...interface{}) {
	fmt.Print(Prefix)
	fmt.Println(v...)
}

func printerr(v ...interface{}) {
	fmt.Fprint(os.Stderr, Prefix)
	fmt.Fprintln(os.Stderr, v...)
}