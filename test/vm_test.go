package gocf

import (
	"fmt"
	"testing"

	"github.com/ryouaki/gocf"
	"github.com/ryouaki/gocf/plugins"
)

func TestVM(t *testing.T) {
	plugins.InitPlugins()
	gocf.InitVM(1)
	vm := gocf.GetVM(1)

	val, err := vm.Ctx.Eval(`(invoke, id) => function () {
		var argvs = [id]
		for (var i = 0; i < arguments.length; i++) {
			var argv = arguments[i];
			argvs.push(argv)
		}
		var ret = invoke.apply(this, argvs);
		try {
			objData = JSON.parse(ret.data)
			ret.data = objData
		} catch(e) {}
		return ret
	}`, "input", 1<<0)

	if val != nil {
		fmt.Print(val.ToString())
	}

	if err != nil {
		e := vm.Ctx.GetException()
		fmt.Print(e.ToString())
	}
}
