package injector

import "fmt"

type (
	Provider[T any] func(*Injector) (T, error)
)

func Provide[T any](i *Injector, provider Provider[T]) error {
	svcName := generateServiceName[T]()
	inj := getInjectorOrDefault(i)
	return inj.setService(svcName, newServiceProvider(svcName, provider))
}

func Invoke[T any](i *Injector) (T, error) {
	svcName := generateServiceName[T]()
	inj := getInjectorOrDefault(i)

	svcAny, ok := inj.getService(svcName)
	if !ok {
		return empty[T](), fmt.Errorf("could not find service `%s`", svcName)
	}
	svc, ok := svcAny.(Service[T])
	if !ok {
		return empty[T](), fmt.Errorf("could not find service `%s`", svcName)
	}
	instance, err := svc.getInstance(i)
	if err != nil {
		return empty[T](), err
	}
	return instance, nil
}
