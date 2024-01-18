package global

import (
	"context"

	"github.com/HrvojeLesar/recommender/config"
	"github.com/HrvojeLesar/recommender/db"
)

type Instance struct {
	Mongo *db.MongoInstance
}

type Global interface {
	context.Context
	Instance() *Instance
}

type GlobalInstances struct {
	context.Context
	instances *Instance
}

func (gi *GlobalInstances) Instance() *Instance {
	return gi.instances
}

func New(ctx context.Context, config config.Config) (Global, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &GlobalInstances{
		Context:   ctx,
		instances: &Instance{},
	}, cancelFunc
}
