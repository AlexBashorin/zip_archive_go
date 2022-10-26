package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	xj "github.com/basgys/goxml2json"
)

func parseXML(xml string) string {
	thisxml := strings.NewReader(xml)
	json, err := xj.Convert(thisxml)
	if err != nil {
		panic("That's embarrassing...")
	}

	return json.String()
}

func main() {
	http.HandleFunc("/parse-xml", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		dec, err := base64.StdEncoding.DecodeString(string(b))

		f, err := os.Create("some.xml")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}

		xmlFile, err := os.Open("./some.xml")
		if err != nil {
			fmt.Println(err)
		}
		defer xmlFile.Close()

		byteValue, _ := ioutil.ReadAll(xmlFile)

		res := parseXML(string(byteValue))

		h, err := json.Marshal(res)
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(h)

		e := os.Remove("./some.xml")
		if e != nil {
			log.Fatal(e)
		}
	})
	log.Fatal(http.ListenAndServe(":6060", nil))
}
