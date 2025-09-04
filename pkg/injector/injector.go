package injector

import (
	"fmt"
	"sync"
)

type (
	Injector struct {
		rwMtx    sync.RWMutex
		services map[string]any
	}
)

var (
	defaultInjector         *Injector
	defaultInjectorInitOnce sync.Once
)

func NewInjector() *Injector {
	return &Injector{services: map[string]any{}}
}

func getInjectorOrDefault(i *Injector) *Injector {
	if i != nil {
		return i
	}
	defaultInjectorInitOnce.Do(func() {
		if defaultInjector == nil {
			defaultInjector = NewInjector()
		}
	})
	return defaultInjector
}

func (inj *Injector) setService(name string, svc any) error {
	inj.rwMtx.Lock()
	defer inj.rwMtx.Unlock()
	_, ok := inj.services[name]
	if ok {
		return fmt.Errorf("service `%s` has already been declared", name)
	}
	inj.services[name] = svc
	return nil
}

func (inj *Injector) getService(name string) (any, bool) {
	inj.rwMtx.RLock()
	defer inj.rwMtx.RUnlock()
	svc, ok := inj.services[name]
	return svc, ok
}
