package client

import (
	"context"

	"github.com/assiljaby/trafic-toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
