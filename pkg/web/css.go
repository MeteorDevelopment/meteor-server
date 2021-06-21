package web

import (
	"bytes"
	"io/ioutil"
	"log"
)

func parseCSS() []byte {
	var buffer bytes.Buffer

	fileBytes, err := ioutil.ReadFile("css/main.css")
	if err != nil {
		log.Fatal(err)
	}

	buffer.Write(fileBytes)

	files, err := ioutil.ReadDir("css")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name() == "main.css" {
			continue
		}

		fileBytes, err = ioutil.ReadFile("css/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		buffer.Write(fileBytes)
	}

	return buffer.Bytes()
}
