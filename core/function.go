package core

/*
#include "./../quickjs-libc.h";
*/
import "C"

type JSGoFuncHandler func(args []*JSValue, this *JSValue) (*JSValue, *JSValue)

type JSGoFunc struct {
	p   C.JSValue
	ctx *JSContext
	fb  JSGoFuncHandler
}
