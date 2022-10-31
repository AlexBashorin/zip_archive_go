package main

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
		// PARSE FORM
		// var buf bytes.Buffer
		errPF := r.ParseMultipartForm(32 << 20)
		// errPF := r.ParseForm()
		check(errPF)

		file, header, errFile := r.FormFile("fileOne")
		check(errFile)
		defer file.Close()

		nameFile := strings.Split(header.Filename, ".")
		ff, errFF := ioutil.ReadFile(header.Filename)
		check(errFF)

		// b, err := io.ReadAll(r.Body)
		// check(err)

		// jsonstr := string(b)
		// var files []Zipit
		// json.Unmarshal([]byte(jsonstr), &files)

		// for i, bod := range files {
		// 	str := bod.Body
		// 	data, err := base64.StdEncoding.DecodeString(str)
		// 	if err != nil {
		// 		log.Fatal("error:", err)
		// 	}
		// 	files[i].Body = string(data)
		// }

		archive, errArch := os.Create("archive.zip")
		check(errArch)
		defer archive.Close()

		wArch := zip.NewWriter(archive)

		f, err := wArch.Create(nameFile[0])
		check(err)
		_, err = f.Write(ff)
		check(err)
		// for _, file := range files {
		// 	f, err := wArch.Create(file.Name)
		// 	check(err)
		// 	_, err = f.Write([]byte(file.Body))
		// 	check(err)
		// }

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
