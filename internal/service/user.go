package service

import (
	"context"
	"effectiveMobileTestTask/internal/entity"
	"effectiveMobileTestTask/internal/store"
	"effectiveMobileTestTask/pkg/age_api"
	"effectiveMobileTestTask/pkg/gender_api"
	"effectiveMobileTestTask/pkg/nationality_api"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type UserActions interface {
	GetUsers(ctx context.Context, params store.UserParamsToFilter) ([]entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	EditUser(ctx context.Context, id int, params store.UserParamsToEdit) error
	AddUser(ctx context.Context, name string, surname string, patronymic interface{}) error
}

type UserService struct {
	log       *logrus.Logger
	userStore store.UserStores
}

func (u *UserService) AddUser(ctx context.Context, name string, surname string, patronymic interface{}) error {
	u.log.Info("[AddUser] started work")

	var (
		wg          sync.WaitGroup
		age         int
		nationality string
		sex         string
		errors      = make(chan error, 3)
	)

	go func() {
		var err error
		defer wg.Done()

		u.log.Info("[AddUser] getting user age")
		age, err = age_api.GetAge(name)
		if err != nil {
			u.log.Errorf("[AddUser] error while getting user age: %v", err.Error())
			errors <- err
			return
		}
		u.log.Info("[AddUser] got user age")
	}()

	go func() {
		var err error
		defer wg.Done()

		u.log.Info("[AddUser] getting user nationality")
		nationality, err = nationality_api.GetNationality(name)
		if err != nil {
			u.log.Errorf("[AddUser] error while getting user nationality: %v", err.Error())
			errors <- err
			return
		}
		u.log.Info("[AddUser] got user nationality")
	}()

	go func() {
		var err error
		defer wg.Done()

		u.log.Info("[AddUser] getting user sex")
		sex, err = gender_api.GetGender(name)
		if err != nil {
			u.log.Errorf("[AddUser] error while getting user sex: %v", err.Error())
			errors <- err
			return
		}
		u.log.Info("[AddUser] got user sex")
	}()

	wg.Add(3)
	wg.Wait()
	close(errors)

	for v := range errors {
		return v
	}

	err := u.userStore.AddUser(ctx, name, surname, patronymic, sex, nationality, age)
	if err != nil {
		u.log.Errorf("[AddUser] error while addind user to db: %s", err.Error())
		return fmt.Errorf("error while adding user to: %s", err.Error())
	}

	u.log.Info("[AddUser] work ended")
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int) error {
	u.log.Info("[DeleteUser] started work")

	err := u.userStore.DeleteUser(ctx, id)
	if err != nil {
		u.log.Errorf("[DeleteUser] error while deleting user from db: %s", err.Error())
		return fmt.Errorf("error while deleting user from db: %s", err.Error())
	}

	u.log.Info("[DeleteUser] work ended")
	return nil
}

func (u *UserService) EditUser(ctx context.Context, id int, params store.UserParamsToEdit) error {
	u.log.Info("[EditUser] started work")

	err := u.userStore.EditUser(ctx, id, params)
	if err != nil {
		u.log.Errorf("[EditUser] error while getting users from db: %s", err.Error())
		return fmt.Errorf("error while getting users from db: %s", err.Error())
	}

	u.log.Info("[EditUser] work ended")
	return nil
}

func (u *UserService) GetUsers(ctx context.Context, params store.UserParamsToFilter) ([]entity.User, error) {
	u.log.Info("[GetAllUsers] started work")

	users, err := u.userStore.GetUsers(ctx, params)
	if err != nil {
		u.log.Errorf("[GetAllUsers] error while getting all users from db: %s", err.Error())
		return nil, fmt.Errorf("error while getting all users from db: %s", err.Error())
	}

	u.log.Info("[GetAllUsers] work ended")
	return users, nil
}

func New(log *logrus.Logger, userStore store.UserStores) *UserService {
	return &UserService{
		log:       log,
		userStore: userStore,
	}
}
