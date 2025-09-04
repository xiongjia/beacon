package injector

import (
	"errors"
	"fmt"
	"sync"
)

type (
	Service[T any] interface {
		getName() string
		getInstance(*Injector) (T, error)
	}

	serviceProvider[T any] struct {
		rwMtx    sync.RWMutex
		name     string
		instance T
		built    bool
		provider Provider[T]
	}
)

func generateServiceName[T any]() string {
	var t T
	name := fmt.Sprintf("%T", t)
	if name != "<nil>" {
		return name
	} else {
		return fmt.Sprintf("%T", new(T))
	}
}

func newServiceProvider[T any](name string, provider Provider[T]) Service[T] {
	return &serviceProvider[T]{
		name:     name,
		built:    false,
		provider: provider,
	}
}

func (p *serviceProvider[T]) getName() string {
	return p.name
}

func (p *serviceProvider[T]) getInstance(i *Injector) (T, error) {
	p.rwMtx.Lock()
	defer p.rwMtx.Unlock()
	if p.built {
		return p.instance, nil
	}
	// build a new instance
	err := p.newServiceInstance(i)
	if err != nil {
		return empty[T](), err
	} else {
		return p.instance, nil
	}
}

func (p *serviceProvider[T]) newServiceInstance(i *Injector) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = errors.New("bad provider")
			}
		}
	}()
	inst, err := p.provider(i)
	if err != nil {
		return err
	}
	p.instance = inst
	p.built = true
	return nil
}
