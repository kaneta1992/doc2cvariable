package main

import (
	"flag"
	"log"

	"github.com/kaneta1992/doc2cvariable/src"
)

func main() {
	outPath := flag.String("o", "shader_code.h", "Set the output filename")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("please input filename")
	}

	d := doc2cvariable.NewDoc2CVariable(*outPath, args)
	err := d.WriteFile()
	if err != nil {
		log.Fatal(err)
	}
}
