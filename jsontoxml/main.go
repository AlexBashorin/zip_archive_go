package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type soft struct {
	Os     string `xml:"os" json:"os"`
	Server string `xml:"server" json:"server"`
	Php    string `xml:"php" json:"php"`
	Db     string `xml:"db" json:"db"`
}

type processor struct {
	Count string `xml:"count,attr" json:"-count"`
	Model string `xml:"model,attr" json:"-model"`
}

type memory struct {
	Memfree  string `xml:"MemFree" json:"MemFree"`
	Memtotal string `xml:"MemTotal" json:"MemTotal"`
}

type volume []struct {
	Mounted    string `xml:"Mounted,attr" json:"-Mounted"`
	Filesystem string `xml:"filesystem,attr" json:"-filesystem"`
	Used       string `xml:"Used,attr" json:"-Used"`
	Available  string `xml:"Available,attr" json:"-Available"`
	Capacity   string `xml:"Capacity,attr" json:"-Capacity"`
}

type disk struct {
	Volume volume `xml:"volume" json:"volume"`
}

type server struct {
	Processor processor
	Memory    memory
	Disk      disk
	Name      string `xml:"name,attr" json:"-name"`
	IP        string `xml:"IP,attr" json:"-IP"`
	Port      string `xml:"port,attr" json:"-port"`
	Soft      soft
}

type xml_doc struct {
	Date   string `xml:"date,attr" json:"-date"`
	Server server
}

func main() {
	startXML := `<?xml version="1.0" encoding="utf-8"?>`

	http.HandleFunc("/parse-json", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}

		jsonStr := string(b)

		var p xml_doc
		json.Unmarshal([]byte(jsonStr), &p)
		out, _ := xml.MarshalIndent(p, "\t", "\t")

		fmt.Println(startXML + string(out))

		w.Header().Set("Content-Type", "text/xml")
		w.Write(out)
	})
	log.Fatal(http.ListenAndServe(":6061", nil))
}
