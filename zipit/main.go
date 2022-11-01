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
		// PARSE FORM
		// // var buf bytes.Buffer
		// errPF := r.ParseMultipartForm(32 << 20)
		// check(errPF)

		// file, handler, errFile := r.FormFile("fileOne")
		// check(errFile)
		// defer file.Close()

		// // f, errOS := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		// fFORM, errOS := os.Create(handler.Filename)
		// check(errOS)
		// defer fFORM.Close()

		// // bodyFILE, err := iotuil.ReadFile("file.txt")
		// // if err != nil {
		// // 	log.Fatalf("unable to read file: %v", err)
		// // }

		// n := "./" + handler.Filename
		// openFILE, _ := os.Open(n)
		// readerFILE := bufio.NewReader(openFILE)
		// contentFILE, _ := ioutil.ReadAll(readerFILE)

		// // nameFile := strings.Split(header.Filename, ".")
		// // ff, errFF := ioutil.ReadFile(header.Filename)
		// // check(errFF)

		file, fileHeader, err := r.FormFile("fileOne")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer file.Close()

		buff := make([]byte, fileHeader.Size)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// // Create the uploads folder if it doesn't
		// // already exist
		// err = os.MkdirAll("./uploads", os.ModePerm)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// // Create a new file in the uploads directory
		// dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// defer dst.Close()

		// // Copy the uploaded file to the filesystem
		// // at the specified destination
		// _, err = io.Copy(dst, file)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		//////// BASE64
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

		fileZip, err := wArch.Create(fileHeader.Filename)
		check(err)
		_, err = fileZip.Write(buff)
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
