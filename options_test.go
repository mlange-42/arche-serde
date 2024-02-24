package archeserde

import (
	"reflect"
	"testing"

	"github.com/mlange-42/arche/generic"
	"github.com/stretchr/testify/assert"
)

type testComp struct{}

func TestOptions(t *testing.T) {
	opt := newSerdeOptions(
		Opts.SkipEntities(),
		Opts.SkipAllComponents(),
		Opts.SkipAllResources(),
		Opts.SkipComponents(generic.T[testComp]()),
		Opts.SkipResources(generic.T[testComp]()),
	)

	assert.True(t, opt.skipEntities)
	assert.True(t, opt.skipAllComponents)
	assert.True(t, opt.skipAllResources)
	assert.Equal(t, []reflect.Type{generic.T[testComp]()}, opt.skipComponents)
	assert.Equal(t, []reflect.Type{generic.T[testComp]()}, opt.skipResources)
}
