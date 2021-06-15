package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/img", getImg)

	fmt.Println("Listening on localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getImg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Accept") == "text/html" {
		w.Header().Set("Content-Type", "text/html")

		htmlFile, err := os.OpenFile("./resources/index.html", os.O_RDWR, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer htmlFile.Close()

		content, err := ioutil.ReadAll(htmlFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		posicion := bytes.Index(content, []byte("img"))

		_, err = htmlFile.Seek(int64(posicion+3), 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		buff := make([]byte, 100)
		nBuff, err := htmlFile.Read(buff)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		var add []byte = []byte(" src='img.png'")

		_, err = htmlFile.WriteAt([]byte(add), int64(posicion+3))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, err = htmlFile.WriteAt(buff[:nBuff], int64(posicion+3+len(add)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		newHtmlFile, err := ioutil.ReadFile("./resources/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		io.WriteString(w, string(newHtmlFile))

		err = ioutil.WriteFile("./resources/index.html", content, 0644)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else {
		w.WriteHeader(http.StatusNotAcceptable)
	}

}
