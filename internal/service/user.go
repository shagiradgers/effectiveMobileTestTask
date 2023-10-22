package service

import (
	"context"
	"effectiveMobileTestTask/internal/entity"
	"effectiveMobileTestTask/internal/store"
	"fmt"
	"github.com/sirupsen/logrus"
)

type UserActions interface {
	GetUsers(ctx context.Context, params store.UserParamsToFilter) ([]*entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	EditUser(ctx context.Context, id int, params store.UserParamsToEdit) error
	AddUser(ctx context.Context, params store.UserParamsToAdd) error
}

type UserService struct {
	log       *logrus.Logger
	userStore UserActions
}

func (u *UserService) AddUser(ctx context.Context, params store.UserParamsToAdd) error {
	u.log.Info("[AddUser] started work")

	err := u.userStore.AddUser(ctx, params)
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

func (u *UserService) GetUsers(ctx context.Context, params store.UserParamsToFilter) ([]*entity.User, error) {
	u.log.Info("[GetAllUsers] started work")

	users, err := u.userStore.GetUsers(ctx, params)
	if err != nil {
		u.log.Errorf("[GetAllUsers] error while getting all users from db: %s", err.Error())
		return nil, fmt.Errorf("error while getting all users from db: %s", err.Error())
	}

	u.log.Info("[GetAllUsers] work ended")
	return users, nil
}

func New(log *logrus.Logger, userStore UserActions) *UserService {
	return &UserService{
		log:       log,
		userStore: userStore,
	}
}
