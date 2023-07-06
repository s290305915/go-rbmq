package rbmq

// 因为go没有枚举类型，所以这里通过常量来模拟枚举
type EnumExchagne int

const (
	DIRECT_EXCHANGE EnumExchagne = iota
	FANOUT_EXCHANGE
	TOPIC_EXCHANGE
	HEADERS_EXCHANGE
)

func (f EnumExchagne) String() string {
	return [...]string{"direct", "fanout", "topic", "headers"}[f]
}
