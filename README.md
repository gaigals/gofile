# GOFile

Package for creating basic files and temporary files for basic testing.


## Usage 

### Struct tags:
Create temporary files with path and given content (cnt).\
Note: Struct field must be exported, have to contain `file` tag and must be FileHandle type.

```go
type MyFiles struct {
    File0      gofile.FileHandle `file:"path:/mypath/hello_world.txt;tmp;cnt:text"`
    File1      gofile.FileHandle `file:"path:/mypath/myfile.txt;tmp;cnt:hello"`
    File2      gofile.FileHandle `file:"tmp;path:/other/test.go"`
}

func main() {
    myFiles := MyFiles{}

    err := NewFileHandleStruct(&myFiles)
    if err != nil {
        log.Fatalln(err)
    }
}
```

### By functions: 

```go
func main() {
NewFileTemp(fullPath, content string, prems os.FileMode)
    myFile, err := gofile.NewFile("/path/to/file", "content", os.ModePerm)
    if err != nil {
        log.Fatalln(err)
    }

    // Creates file in temporary directory, handled by OS/GOLang.
    myTempFile, err := gofile.NewFileTemp("/path/to/file", "content", os.ModePerm)
    if err != nil {
        log.Fatalln(err)
    }
}
```

### Writing and reading 
```go
type MyFiles struct {
    File gofile.FileHandle `file:"path:/dir/file";tmp`
}

func main() {
    myFiles := MyFiles{}
    
    err := NewFileHandleStruct(&myFiles)
    if err != nil {
        log.Fatalln(err)
    }
    
    content := "some random text"
    
    err = myFiles.File.Write(content, os.ModePrem)
    if err != nil {
        log.Fatalln(err)
    }
    
    fileContent, err := myFiles.File.Read()
    if err != nil {
        log.Fatalln(err)
    }
    
    err = myFiles.File.Remove()
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(fileContent)
}
```
