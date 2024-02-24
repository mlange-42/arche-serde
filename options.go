package archeserde

import (
	"reflect"

	"github.com/mlange-42/arche/generic"
)

var Opts = Options{}

type Option func(o *serdeOptions)

type Options struct{}

func (o Options) SkipAllResources() Option {
	return func(o *serdeOptions) {
		o.skipAllResources = true
	}
}

func (o Options) SkipAllComponents() Option {
	return func(o *serdeOptions) {
		o.skipAllComponents = true
	}
}

func (o Options) SkipEntities() Option {
	return func(o *serdeOptions) {
		o.skipEntities = true
	}
}

func (o Options) SkipComponents(comps ...generic.Comp) Option {
	return func(o *serdeOptions) {
		o.skipComponents = make([]reflect.Type, len(comps))
		for i, c := range comps {
			o.skipComponents[i] = reflect.Type(c)
		}
	}
}

func (o Options) SkipResources(comps ...generic.Comp) Option {
	return func(o *serdeOptions) {
		o.skipResources = make([]reflect.Type, len(comps))
		for i, c := range comps {
			o.skipResources[i] = reflect.Type(c)
		}
	}
}

type serdeOptions struct {
	skipAllResources  bool
	skipAllComponents bool
	skipEntities      bool

	skipComponents []reflect.Type
	skipResources  []reflect.Type
}

func newSerdeOptions(opts ...Option) serdeOptions {
	o := serdeOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
