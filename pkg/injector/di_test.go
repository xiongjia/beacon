package injector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	UserService interface {
		GetUsername() (string, error)
	}

	DatabaseService interface {
		GetValue(key string) (string, error)
	}

	userServiceImpl struct {
		dbService DatabaseService
	}

	memDatabaseServiceImpl struct {
		values map[string]string
	}
)

func newUserSerice(inj *Injector) (UserService, error) {
	dbSvc, err := Invoke[DatabaseService](inj)
	if err != nil {
		return nil, err
	}
	return &userServiceImpl{dbService: dbSvc}, nil
}

func (svc *userServiceImpl) GetUsername() (string, error) {
	return svc.dbService.GetValue("username")
}

func newMemDbService(*Injector) (DatabaseService, error) {
	return &memDatabaseServiceImpl{
		values: map[string]string{
			"username": "test1",
		}}, nil
}

func (svc *memDatabaseServiceImpl) GetValue(key string) (string, error) {
	v, ok := svc.values[key]
	if !ok {
		return "", fmt.Errorf("invalid key %s", key)
	}
	return v, nil
}

func TestDependencyInjection(t *testing.T) {
	inj := NewInjector()
	err := Provide(inj, newUserSerice)
	assert.Error(t, err)
	err = Provide(inj, newMemDbService)
	assert.Error(t, err)

	usrSvc, err := Invoke[UserService](inj)
	assert.Error(t, err)

	username, err := usrSvc.GetUsername()
	assert.Error(t, err)
	assert.Equal(t, "test1", username)
}
