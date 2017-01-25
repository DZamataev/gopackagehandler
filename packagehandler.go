package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var globalPath = "./"

func main() {
	pathPtr := flag.String("path", "./", "directory to save uploaded files")
	portPtr := flag.Int("port", 8085, "port to listen")

	flag.Parse()
	fmt.Println("path:", *pathPtr)
	fmt.Println("port:", *portPtr)

	globalPath = *pathPtr

	http.HandleFunc("/upload", upload)
	http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), nil)
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	// curl -i -X POST -H "Content-Type: multipart/form-data" -F "data=@mv.deb" http://localhost:8080/upload/

	fmt.Println("request method:", r.Method)
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		f, err := os.Create(globalPath + handler.Filename)
		if err != nil {
			fmt.Println("Error creating file")
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		fmt.Println("File uploaded.")

		cmd := exec.Command("tr", "a-z", "A-Z") // reprepro -b /var/packages includedeb xenial example-helloworld_1.0.0.0_*
		cmd.Stdin = strings.NewReader("some input")
		var out bytes.Buffer
		cmd.Stdout = &out
		execError := cmd.Run()
		if execError != nil {
			fmt.Println("Error executing reprepro command")
			fmt.Println(execError)
		} else {
			fmt.Println("Reprepro output:")
			fmt.Printf("%q\n", out.String())
		}
	}
}
