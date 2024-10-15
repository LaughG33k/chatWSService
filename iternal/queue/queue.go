package queue

import "context"

type Publisher interface {
	Publish(ctx context.Context, receiver string, message any) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, receiver string, ch chan any) error
}

type Queue interface {
	Publisher
	Subscriber
}
