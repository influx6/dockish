package docker

import (
	"context"

	"github.com/influx6/faux/context"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/ops"
	"github.com/moby/moby/client"
)

// ImageTag returns a new ImageTagOp instance to be executed on the client.
func (d *DockerCaster) ImageTag(tag string) (*ImageTagOp, error) {
	var spell ImageTagOp

	spell.tag = tag

	return &spell, nil
}

// ImageTagOptions defines a function type to modify internal fields of the ImageTagOp.
type ImageTagOptions func(*ImageTagOp)

// ImageTagResponseCallback defines a function type for ImageTagOp response.
type ImageTagResponseCallback func() error

// ImageTagOp defines a structure which implements the Op interface
// for executing of docker based commands for ImageTag.
type ImageTagOp struct {
	client *client.Client

	tag string
}

// Op returns a object implementing the ops.Op interface.
func (cm *ImageTagOp) Op(callback ImageTagResponseCallback) ops.Op {
	return &onceImageTagOp{spell: cm, callback: cb}
}

type onceImageTagOp struct {
	callback ImageTagResponseCallback
	spell    *ImageTagOp
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageTagOp) Exec(ctx context.CancelContext, m metrics.Metrics) error {
	return cm.spell.Exec(ctx, m, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageTagOp) Exec(ctx context.CancelContext, m metrics.Metrics, callback ImageTagResponseCallback) error {
	if cm.client == nil {
		return ErrNoDockerClientProvided
	}

	done := make(chan struct{})
	defer close(done)

	// Cancel context if are done or if context has expired.
	reqCtx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			cancel()
			return
		case <-done:
			return
		}
	}()

	// Execute client ImageTag method.
	err := cm.client.ImageTag(reqCtx, cm.tag)
	if err != nil {
		return err
	}

	if callback != nil {
		return callback()
	}

	return nil
}
