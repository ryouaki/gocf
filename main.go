package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ryouaki/gocf/core"
)

func init() {

}

func main() {
	f, err := os.OpenFile("./demo/main.js", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Read Script Failed:", err)
		f.Close()
		return
	}

	data, err1 := ioutil.ReadAll(f)
	if err1 != nil {
		fmt.Println("Read Script Failed:", err1)
	}

	core.InitGoCloudFunc(string(data))
}
