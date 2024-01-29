package rbmq

import (
	"context"
	"fmt"
	"unsafe"
)

type iface struct {
	itab, data uintptr
}

type emptyCtx int

type valueCtx struct {
	context.Context
	key, val any
}

func GetKeyValues(ctx context.Context) map[string]any {
	m := make(map[string]any)
	getKeyValue(ctx, m)
	return m
}

func getKeyValue(ctx context.Context, m map[string]any) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 || int(*(*emptyCtx)(unsafe.Pointer(ictx.data))) == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil {
		key := fmt.Sprintf("%v", valCtx.key)
		m[key] = valCtx.val
	}
	getKeyValue(valCtx.Context, m)
}
