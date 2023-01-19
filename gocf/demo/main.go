package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ryouaki/gocf"
	"github.com/ryouaki/gocf/plugins"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "../quickjs-libc.h"
*/
import "C"

func init() {

	plugins.InitPlugins()

	gocf.InitGoCloudFunc()
}

func main() {
	f, err := os.OpenFile("./main.js", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Read Script Failed:", err)
		f.Close()
		return
	}

	src, err1 := ioutil.ReadAll(f)
	if err1 != nil {
		fmt.Println("Read Script Failed:", err1)
	}

	gocf.RunAPI(string(src))
}
