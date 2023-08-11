package main

import (
	"github.com/gaigals/gofile"
	"log"
	"os"
)

type MyFiles struct {
	File gofile.FileHandle `file:"path:/dir/file;tmp"`
}

func main() {
	myFiles := MyFiles{}

	err := gofile.NewFileHandleStruct(&myFiles)
	if err != nil {
		log.Fatalln(err)
	}

	content := "some random text"

	err = myFiles.File.Write(content, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}
