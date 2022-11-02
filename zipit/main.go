package main

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type Zipit struct {
	Name, Body string
}

func main() {
	// write zip
	http.HandleFunc("/zipit", func(w http.ResponseWriter, r *http.Request) {
		archive, errArch := os.Create("archive.zip")
		check(errArch)
		defer archive.Close()
		wArch := zip.NewWriter(archive)

		err := r.ParseMultipartForm(32 << 20)
		check(err)

		fhs := r.MultipartForm.File["fileOne"]
		for _, fh := range fhs {
			f, err := fh.Open()
			check(err)

			fName, errName := wArch.Create(fh.Filename)
			check(errName)

			fBuff := make([]byte, fh.Size)
			_, errf := f.Read(fBuff)
			if errf != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_, errWrite := fName.Write(fBuff)
			check(errWrite)
		}

		errWrite := wArch.Close()
		check(errWrite)

		a, _ := os.Open("./archive.zip")
		reader := bufio.NewReader(a)
		content, _ := ioutil.ReadAll(reader)

		encoded := base64.StdEncoding.EncodeToString(content)

		w.Header().Set("Content-Type", "text/json")
		w.Write([]byte(encoded))

		e := os.Remove("./archive.zip")
		check(e)
	})
	log.Fatal(http.ListenAndServe(":5050", nil))
}
