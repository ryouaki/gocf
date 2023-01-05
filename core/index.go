package core

/*
#include "./../quickjs-libc.h";
#include "./invoke.h"
*/
import "C"

func InitGoCloudFunc() {
	C.Invoke(nil, C.JSValue{}, C.int(0), nil)
}
