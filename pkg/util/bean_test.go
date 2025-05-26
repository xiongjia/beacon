package util

import (
	"fmt"
	"testing"
)

type (
	UserRepository struct{}

	UserService interface {
		GetUser(id int) string
	}

	UserServiceImpl struct {
		Repo *UserRepository `inject:"userRepository"`
	}
)

func (s *UserServiceImpl) AfterPropertiesSet() {
	fmt.Println("UserService initialized")
}

func (s *UserServiceImpl) GetUser(id int) string {
	return s.Repo.Find(id)
}

func (r *UserRepository) Find(id int) string {
	return fmt.Sprintf("User-%d", id)
}

func TestBeanContaner(t *testing.T) {
	container := NewBeanContainer()
	container.RegisterBean("userRepository", func() *UserRepository {
		return &UserRepository{}
	}, BeanScopeSingleton)

	container.RegisterBean("userService", func() UserService {
		return &UserServiceImpl{}
	}, BeanScopeSingleton)

	userService, err := container.GetBean("userService")
	if err != nil {
		return
	}
	svc := userService.(UserService)
	userId := svc.GetUser(123)
	t.Logf("user id = %s", userId)
}
