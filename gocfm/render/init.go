package render

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

var base string = ""

func init() {
	if os.Getenv("MODE") == "DEV" {
		base = "../"
	}
}

func buildTemplate(entry string, data any, temps ...string) ([]byte, error) {
	var tempData []byte
	tpl, err := template.New(entry).ParseFiles(temps...)
	if err != nil {
		errString := fmt.Sprintf("System error %s", err)
		tempData = []byte(errString)
		return tempData, err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)

	if err != nil {
		errString := fmt.Sprintf("System error %s", err)
		tempData = []byte(errString)
		return tempData, err
	}

	tempData = buf.Bytes()
	return tempData, nil
}
