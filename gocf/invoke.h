#include "quickjs-libc.h"

JSValue Invoke(JSContext *ctx, JSValueConst thisCtx, int argc, JSValueConst *argv);
JSContext* NewJsContext(JSRuntime *rt);
JSModuleDef* GetModule(JSValueConst importVal);
JSContext* ListModule(JSContext *ctx);