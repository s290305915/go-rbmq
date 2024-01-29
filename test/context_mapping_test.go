package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/s290305915/go-rbmq/rbmq"
)

type TestData struct {
	Name   string
	Number int
}

func TestContextToMapping(t *testing.T) {
	// build a context for test
	ctx := context.Background()

	ctx = context.WithValue(ctx, "key1", "value1")

	ctx, _ = context.WithCancel(ctx)

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key2",
		Number: 1,
	}, "value2")

	ctx, _ = context.WithDeadline(ctx, time.Now().Add(time.Second*20))

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key3",
		Number: 2,
	}, "value3")

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key4",
		Number: 3,
	}, nil)

	ctx = context.WithValue(ctx, &TestData{
		Name:   "key5",
		Number: 4,
	}, nil)

	ctx = context.WithValue(ctx, "ProductLine", "PD1")
	ctx = context.WithValue(ctx, "SaasToken", "")
	ctx = context.WithValue(ctx, "AppId", "PD1_APP")

	ctx = context.WithValue(ctx, "noise_key", "123")
	ctx = context.WithValue(ctx, "ProductLine_extra", map[string]any{"gg1": "123", "gg2": 99, "gg3": struct {
		aa string
		bb int
	}{
		aa: "123",
		bb: 99,
	}})

	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	ctx = context.WithValue(ctx, "ddd", "123")

	// get keys and values
	m := rbmq.GetKeyValues(ctx)
	printMapKeyValue(m)

	// output:
	// [key: &{Name:key3 Number:2}] [value: value3]
	// [key: &{Name:key2 Number:1}] [value: value2]
	// [key: key1] [value: value1]
	// [key: &{Name:key5 Number:4}] [value: <nil>]
	// [key: &{Name:key4 Number:3}] [value: <nil>]
}

func printMapKeyValue(m []rbmq.CtxData) {
	for _, v := range m {
		fmt.Printf("[key: %s value: %+v]\n", v.Key, v.Value)
	}
}
