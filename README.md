# go-cloud-function

1. 用户通过mapi接口管理接口
- 1. mapi/check 接口，监控
- 2. mapi/scripts 接口，提交脚本进入开发环境。可以测试脚本是否正确
- 3. mapi/release 接口，将开发脚本同步到线上环境

2. 用户通过提交js文件自动创建api接口
- 1. api目录开头的文件会自动创建接口
- 2. api目录下，路径名+文件名等于接口地址
- 3. api目录下文件名必须以get，post等method允许的名字开头，用" . "进行分割
- 4. 如果不符合规则将自动忽略。


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

<link data-name="vs/editor/editor.main" rel="stylesheet" href="node_modules/monaco-editor/min/vs/editor/editor.main.css">
<script>
  var require = { paths: { vs: 'node_modules/monaco-editor/min/vs' } };
</script>
<script src="node_modules/monaco-editor/min/vs/loader.js"></script>
<script src="node_modules/monaco-editor/min/vs/editor/editor.main.nls.js"></script>
<script src="node_modules/monaco-editor/min/vs/editor/editor.main.js"></script>