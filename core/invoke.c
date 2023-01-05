#include "_cgo_export.h";

JSValue Invoke(JSContext *ctx, JSValueConst thisCtx, int argc, JSValueConst *argv){
  return GoInvoke(ctx, thisCtx, argc, argv);
}
