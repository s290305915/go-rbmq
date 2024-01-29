package rbmq

import (
	"context"
	"fmt"
	"unsafe"
)

type iface struct {
	itab, data uintptr
}

type CtxData struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type emptyCtx int

type valueCtx struct {
	context.Context
	key, val any
}

func GetKeyValues(ctx context.Context) []CtxData {
	m := make([]CtxData, 0)
	getKeyValue(ctx, &m)
	return m
}

func getKeyValue(ctx context.Context, m *[]CtxData) {
	ictx := *(*iface)(unsafe.Pointer(&ctx))
	if ictx.data == 0 || int(*(*emptyCtx)(unsafe.Pointer(ictx.data))) == 0 {
		return
	}

	valCtx := (*valueCtx)(unsafe.Pointer(ictx.data))
	if valCtx != nil && valCtx.key != nil {
		key := fmt.Sprintf("%+v", valCtx.key)
		*m = append(*m, CtxData{
			Key:   key,
			Value: valCtx.val,
		})
	}
	getKeyValue(valCtx.Context, m)
}
