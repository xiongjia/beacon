package injector

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	UserService interface {
		GetUsername() (string, error)
		GetRandNum() int32
	}

	DatabaseService interface {
		GetValue(key string) (string, error)
		GetRandNum() int32
	}

	userServiceImpl struct {
		dbService DatabaseService
		randnum   int32
	}

	memDatabaseServiceImpl struct {
		values  map[string]string
		randnum int32
	}
)

func newUserSerice(inj *Injector) (UserService, error) {
	dbSvc, err := Invoke[DatabaseService](inj)
	if err != nil {
		return nil, err
	}
	return &userServiceImpl{dbService: dbSvc, randnum: rand.Int32()}, nil
}

func (svc *userServiceImpl) GetUsername() (string, error) {
	return svc.dbService.GetValue("username")
}

func (svc *userServiceImpl) GetRandNum() int32 {
	return svc.randnum
}

func newMemDbService(*Injector) (DatabaseService, error) {
	return &memDatabaseServiceImpl{
		values: map[string]string{
			"username": "test1",
		},
		randnum: rand.Int32(),
	}, nil
}

func (svc *memDatabaseServiceImpl) GetValue(key string) (string, error) {
	v, ok := svc.values[key]
	if !ok {
		return "", fmt.Errorf("invalid key %s", key)
	}
	return v, nil
}

func (svc *memDatabaseServiceImpl) GetRandNum() int32 {
	return svc.randnum
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

	// usrSvc & usrSvc2 should refer to the same service instance
	usrSvc2, err := Invoke[UserService](inj)
	assert.Error(t, err)
	assert.Equal(t, usrSvc.GetRandNum(), usrSvc2.GetRandNum())
}
