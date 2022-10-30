package main

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// write zip
	type Zipit struct {
		Name, Body string
	}

	http.HandleFunc("/zipit", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		jsonstr := string(b)
		var files []Zipit
		json.Unmarshal([]byte(jsonstr), &files)

		for i, bod := range files {
			str := bod.Body
			data, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				log.Fatal("error:", err)
			}
			files[i].Body = string(data)
		}

		archive, errArch := os.Create("archive.zip")
		if errArch != nil {
			panic(errArch)
		}
		defer archive.Close()

		wArch := zip.NewWriter(archive)

		for _, file := range files {
			f, err := wArch.Create(file.Name)
			if err != nil {
				log.Fatal(err)
			}
			_, err = f.Write([]byte(file.Body))
			if err != nil {
				log.Fatal(err)
			}
		}

		errWrite := wArch.Close()
		if errWrite != nil {
			log.Fatal(errWrite)
		}

		a, _ := os.Open("./archive.zip")
		reader := bufio.NewReader(a)
		content, _ := ioutil.ReadAll(reader)

		encoded := base64.StdEncoding.EncodeToString(content)

		w.Header().Set("Content-Type", "text/json")
		w.Write([]byte(encoded))

		e := os.Remove("./archive.zip")
		if e != nil {
			log.Fatal(e)
		}
	})
	log.Fatal(http.ListenAndServe(":5050", nil))
}
