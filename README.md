# go-cloud-function

1. 用户创建接口。并存入js代码【草稿，开发模式，测试模式，上线模式】
2. 用户创建的接口需要提供运行按钮以测试是否可以在gocf正常运行
3. 提供插件模块，插件模块可以注入各种插件供JS调用go模块
4. 提供http请求插件

[] JS调C
[] JS调GO
[] GO调JS
[] GO调C


#cgo CFLAGS: -I./3rdparty/include/quickjs
#cgo linux,!android,386 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/linux/x86 -lquickjs
#cgo linux,!android,amd64 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/linux/x86_64 -lquickjs
#cgo linux,!android LDFLAGS: -lm -ldl -lpthread
#cgo windows,386 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/windows/x86 -lquickjs
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/windows/x86_64 -lquickjs
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/darwin -lquickjs
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/darwin/arm64 -lquickjs
#cgo android,386 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/Android/x86 -lquickjs
#cgo android,amd64 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/Android/x86_64 -lquickjs
#cgo android,arm LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/Android/armeabi-v7a -lquickjs
#cgo android,arm64 LDFLAGS: -L${SRCDIR}/3rdparty/libs/quickjs/Android/arm64-v8a -lquickjs
#cgo android LDFLAGS: -landroid -llog -lm