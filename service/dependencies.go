package service

import (
	"context"
	"github.com/golang/protobuf/proto"
)

type (
	Publisher interface {
		Publish(ctx context.Context, routingKey string, m proto.Message) error
	}
)
