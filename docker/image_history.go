package docker

import (
	"context"

	"github.com/docker/docker/api/types/image"
	"github.com/influx6/box"
	"github.com/influx6/faux/context"
	"github.com/moby/moby/client"
)

// ImageHistory returns a new ImageHistorySpell instance to be executed on the client.
func (d *DockerCaster) ImageHistory() (*ImageHistorySpell, error) {
	var spell ImageHistorySpell

	return &spell, nil
}

// ImageHistorySpell defines a function type to modify internal fields of the ImageHistorySpell.
type ImageHistoryOptions func(*ImageHistorySpell)

// ImageHistoryResponseCallback defines a function type for ImageHistorySpell response.
type ImageHistoryResponseCallback func(image.HistoryResponseItem) error

// ImageHistorySpell defines a structure which implements the Spell interface
// for executing of docker based commands for ImageHistory.
type ImageHistorySpell struct {
	client *client.Client
}

// Spell returns a object implementing the box.Shell interface.
func (cm *ImageHistorySpell) Spell(callback ImageHistoryResponseCallback) box.Spell {
	return &onceImageHistorySpell{spell: cm, callback: cb}
}

type onceImageHistorySpell struct {
	callback ImageHistoryResponseCallback
	spell    *ImageHistorySpell
}

// Exec excutes the spell and adds the neccessary callback.
func (cm *onceImageHistorySpell) Exec(ctx context.CancelContext) error {
	return cm.spell.Exec(ctx, cm.callback)
}

// Exec executes the image creation for the underline docker server pointed to.
func (cm *ImageHistorySpell) Exec(ctx context.CancelContext, callback ImageHistoryResponseCallback) error {
	// Execute client ImageHistory method.
	ret0, err := cm.client.ImageHistory()
	if err != nil {
		return err
	}

	if callback != nil {
		return callback(ret0)
	}

	return nil
}
