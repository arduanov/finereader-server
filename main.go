package main

import (
	"log"
	"github.com/go-gem/gem"
	"os"
	"io"
	"path/filepath"
	"flag"
)

var ocrExec = flag.String("ocr", `C:\Program Files\ABBYY FineReader 12\FineCmd.exe`, "OCR exec path")
var inPath = flag.String("in", "./in/", "in path")
var outPath = flag.String("out", "./out/", "out path")
var in = ""
var out = ""

func main() {
	// Create dirs
	in, _ = filepath.Abs(*inPath)
	out, _ = filepath.Abs(*outPath)
	os.MkdirAll(in, 0755)
	os.MkdirAll(out, 0755)
	// Create server.
	srv := gem.New(":8080")

	// Create router.
	router := gem.NewRouter()
	// Register handler
	router.POST("/", index)

	// Start server.
	log.Println(srv.ListenAndServe(router.Handler()))
}

func index(ctx *gem.Context) {
	w := ctx.Response

	filename, err := handleUpload(ctx, in)
	if err != nil {
		panic(err) //(err)
		//log.Println(err)
		return
	}
	//ocr
	inFile := filepath.Join(in, filename)
	outFile := filepath.Join(out, filename)

	err = ocr(*ocrExec, inFile, outFile)
	if err != nil {
		panic(err) //(err)
		//log.Println(err)
		return
	}
	//send file
	fh, err := os.Open(filepath.Join(out, filename))
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	_, err = io.Copy(w, fh)
	if err != nil {
		panic(err) //(err)
		//log.Println(err)
		return
	}
	//remove source and dest
	os.Remove(filepath.Join(in, filename))
	os.Remove(filepath.Join(out, filename))
}

func handleUpload(ctx *gem.Context, dstPath string) (filename string, err error) {
	file, handler, err := ctx.FormFile("file")
	if err != nil {
		panic(err) //(err)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(filepath.Join(dstPath, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err) //(err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		panic(err) //(err)
		return
	}
	err = f.Sync()
	if err != nil {
		panic(err) //(err)
		return
	}
	return handler.Filename, err
}
