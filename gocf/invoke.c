#include "_cgo_export.h"

JSValue Invoke(JSContext *ctx, JSValueConst thisCtx, int argc, JSValueConst *argv){
  return GoInvoke(ctx, thisCtx, argc, argv);
}

JSContext* NewJsContext(JSRuntime *rt) {
	JSContext* ctx = JS_NewContext(rt);
	js_std_init_handlers(rt);
	JS_SetModuleLoaderFunc(rt, NULL, js_module_loader, NULL);

	return ctx;
}

JSModuleDef* GetModule(JSValueConst importVal) {
  return JS_VALUE_GET_PTR(importVal);
}

JSContext* ListModule(JSContext *ctx) {
	// ctx->loaded_modules;

	return ctx;
}