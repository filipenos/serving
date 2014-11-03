package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	port, down, up string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&port, "port", "8080", "Serve port number")
	flag.StringVar(&down, "download-dir", pwd, "Directory to be served")
	flag.StringVar(&up, "upload-dir", pwd, "Directory to upload files")
}

func main() {
	flag.Parse()
	fmt.Println("Server start on: ", port)
	fmt.Println("Directory to serve: ", down)
	fmt.Println("Directory to upload files: ", up)

	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir(down)))

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Failed to start server, %v\n", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "POST":
		err := r.ParseMultipartForm(32 << 20) // 32MB is the default size used by FormFile
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := r.MultipartForm
		files := m.File["files"]
		for i := range files {
			file, err := files[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			dst, err := os.Create(fmt.Sprintf("%s/%s", up, files[i].Filename))
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			display(w, nil)
		}

	case "GET":
		display(w, nil)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func display(w http.ResponseWriter, data interface{}) {
	t, err := template.ParseFiles("upload.gtpl")
	if err != nil {
		panic(fmt.Sprintf("An error ocurred when parsing template, %v\n", err))
	}
	t.Execute(w, nil)
}
