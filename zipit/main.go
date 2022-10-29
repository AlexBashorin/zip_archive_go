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
	// read zip
	// r, err := zip.OpenReader("./test_zip.zip")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer r.Close()

	// for _, f := range r.File {
	// 	fmt.Printf("Contents of %s: \n", f.Name)
	// 	rc, err := f.Open()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	_, err = io.CopyN(os.Stdout, rc, 68)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	rc.Close()
	// }

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

		archive, errArch := os.Create("archive.zip")
		if errArch != nil {
			panic(errArch)
		}
		defer archive.Close()

		// buf := new(bytes.Buffer)
		wArch := zip.NewWriter(archive)

		// var files = []struct {
		// 	Name, Body string
		// }{
		// 	{"readme.txt", "This archive contains some text files."},
		// 	{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		// 	{"todo.txt", "Get animal handling licence.\nWrite more examples."},
		// }

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
