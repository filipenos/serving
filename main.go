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
		fmt.Printf("Failed to start server, %v", err)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Download"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "POST":
		err := r.ParseMultipartForm(100000)
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
			//http.Redirect(w, r, "/upload", 200)
			display(w, nil)
		}

	case "GET":
		display(w, nil)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func display(w http.ResponseWriter, data interface{}) {
	t := template.New("upload")
	t, _ = t.Parse(upload)
	t.ExecuteTemplate(w, "upload", data)
}
