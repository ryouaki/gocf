#include "_cgo_export.h"

#include "stdio.h"

JSValue Invoke(JSContext *ctx, JSValueConst thisCtx, int argc, JSValueConst *argv){
  return GoInvoke(ctx, thisCtx, argc, argv);
}

JSContext* NewJsContext(JSRuntime *rt) {
	JSContext* ctx = JS_NewContext(rt);
	js_std_init_handlers(rt);
	JS_SetModuleLoaderFunc(rt, NULL, js_module_loader, NULL);

	return ctx;
}

JSModuleDef* GetModule(JSContext *ctx, JSAtom module_name) {
	JSModuleDef *m;
	m = JS_FindLoadedModule(ctx, module_name);
	if (m) {
			JS_FreeAtom(ctx, module_name);
			return m;
	}

	return m;
}

void FreeModule(JSContext *ctx, JSModuleDef *m) {
	return JS_FreeModule(ctx, m);
}