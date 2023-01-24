#include "quickjs-libc.h"

JSValue Invoke(JSContext *ctx, JSValueConst thisCtx, int argc, JSValueConst *argv);
JSContext* NewJsContext(JSRuntime *rt);
JSModuleDef* GetModule(JSContext *ctx, JSAtom module_name);
void FreeModule(JSContext *ctx, JSModuleDef *m);