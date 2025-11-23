package eventbus

type ProduceOpt func(option *ProduceOption)

type ProduceOption struct {
	ShardingKey *string
}

func WithShardingKey(key string) ProduceOpt {
	return func(o *ProduceOption) {
		o.ShardingKey = &key
	}
}
