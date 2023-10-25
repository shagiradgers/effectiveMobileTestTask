package store

import (
	"context"
	"database/sql"
	"effectiveMobileTestTask/internal/entity"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"reflect"
	"time"
)

const (
	CtxTimeout = 1 * time.Second
)

type UserStores interface {
	GetUsers(ctx context.Context, params UserParamsToFilter) ([]entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	EditUser(ctx context.Context, id int, params UserParamsToEdit) error
	AddUser(ctx context.Context, name string, surname string, patronymic interface{}, sex interface{}, nationality interface{}, age interface{}) error
}

// UserParamsToFilter - Параметры для фильтрации пользователя
type UserParamsToFilter struct {
	Id          interface{}
	Age         interface{}
	Name        interface{}
	Surname     interface{}
	Patronymic  interface{}
	Nationality interface{}
	Sex         interface{}
	Limit       interface{}
	Page        interface{}
}

type UserParamsToEdit struct {
	Name        interface{}
	Surname     interface{}
	Patronymic  interface{}
	Sex         interface{}
	Nationality interface{}
	Age         interface{}
}

type InvalidParam struct {
	value string
}

func (i *InvalidParam) Error() string {
	return fmt.Sprintf("Invalid param %s", i.value)
}

type UserStore struct {
	db *sqlx.DB
}

func (u *UserStore) GetUsers(ctx context.Context, params UserParamsToFilter) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	sqlBuilder := squirrel.
		Select("id, name, surname, patronymic, age, sex, nationality").
		From("users")

	if params.Id != 0 {
		sqlBuilder = sqlBuilder.Where("id = ?", params.Id)
	}
	if v, ok := params.Name.(string); ok {
		sqlBuilder = sqlBuilder.Where("name = ?", v)
	}
	if v, ok := params.Surname.(string); ok {
		sqlBuilder = sqlBuilder.Where("surname = ?", v)
	}
	if v, ok := params.Patronymic.(string); ok {
		sqlBuilder = sqlBuilder.Where("patronymic = ?", v)
	}
	if v, ok := params.Nationality.(string); ok {
		sqlBuilder = sqlBuilder.Where("nationality = ?", v)
	}
	if v, ok := params.Sex.(string); ok {
		sqlBuilder = sqlBuilder.Where("sex = ?", v)
	}
	if v, ok := params.Age.(int); ok {
		sqlBuilder = sqlBuilder.Where("age = ?", v)
	}
	if limit, ok := params.Limit.(uint64); ok {
		if page, ok := params.Page.(uint64); ok {
			sqlBuilder = sqlBuilder.Offset(page * limit)
		}
		sqlBuilder = sqlBuilder.Limit(limit)
	}

	query, args, err := sqlBuilder.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	usersDAO := make([]entity.UserDAO, 0)
	err = u.db.SelectContext(ctx, &usersDAO, query, args...)
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, len(usersDAO))
	for i, user := range usersDAO {
		users[i] = entity.User{
			Id:          user.Id,
			Name:        user.Name,
			Surname:     user.Surname,
			Patronymic:  user.Patronymic.String,
			Age:         user.Age,
			Nationality: user.Nationality,
			Sex:         user.Sex,
		}
	}

	return users, nil
}

func (u *UserStore) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	sqlBuilder := squirrel.
		Delete("users").
		Where("id = ?", id)

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) EditUser(ctx context.Context, id int, params UserParamsToEdit) error {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	sqlBuilder := squirrel.
		Update("users").
		Where("id = ?", id)

	if v, ok := params.Name.(string); ok {
		sqlBuilder = sqlBuilder.Set("name", v)
	}
	if v, ok := params.Surname.(string); ok {
		sqlBuilder = sqlBuilder.Set("surname", v)
	}
	if v, ok := params.Patronymic.(string); ok {
		sqlBuilder = sqlBuilder.Set("patronymic", v)
	}
	if v, ok := params.Nationality.(string); ok {
		sqlBuilder = sqlBuilder.Set("nationality", v)
	}
	if v, ok := params.Sex.(string); ok {
		sqlBuilder = sqlBuilder.Set("sex", v)
	}
	if v, ok := params.Age.(int); ok {
		sqlBuilder = sqlBuilder.Set("age", v)
	}

	query, args, err := sqlBuilder.
		ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func getStringInNullStringFromInterface(toStr interface{}) (sql.NullString, error) {
	if toStr == nil {
		return sql.NullString{}, nil
	}
	if v, ok := toStr.(string); !ok {
		return sql.NullString{}, &InvalidParam{value: reflect.TypeOf(toStr).Name()}
	} else {
		if v == "" {
			return sql.NullString{Valid: false}, nil
		}
		return sql.NullString{String: v, Valid: true}, nil
	}
}

func getIntInNullStringFromInterface(toInt interface{}) (sql.NullInt64, error) {
	if toInt == nil {
		return sql.NullInt64{}, nil
	}
	if v, ok := toInt.(int); !ok {
		return sql.NullInt64{}, &InvalidParam{value: reflect.TypeOf(toInt).Name()}
	} else {
		if v == 0 {
			return sql.NullInt64{Valid: false}, nil
		}
		return sql.NullInt64{Int64: int64(v), Valid: true}, nil
	}
}

// AddUser
// ctx - context
// name - username
// surname - user surname
// patronymic - user patronymic (string) (can be nil)
// sex - user sex (string) (can be nil)
// nationality - user nationality (string) (can be nil)
// age - user age (int) (can be nil)
func (u *UserStore) AddUser(ctx context.Context, name string, surname string, patronymic interface{}, sex interface{}, nationality interface{}, age interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	values := []interface{}{name, surname}

	sqlBuilder := squirrel.
		Insert("users").
		Columns("name, surname, patronymic, nationality, sex, age")

	if v, err := getStringInNullStringFromInterface(patronymic); err != nil {
		return err
	} else {
		values = append(values, v)
	}
	if v, err := getStringInNullStringFromInterface(sex); err != nil {
		return err
	} else {
		values = append(values, v)
	}
	if v, err := getStringInNullStringFromInterface(nationality); err != nil {
		return err
	} else {
		values = append(values, v)
	}
	if v, err := getIntInNullStringFromInterface(age); err != nil {
		return err
	} else {
		values = append(values, v)
	}

	query, args, err := sqlBuilder.
		Values(values...).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func NewUserStore(db *sqlx.DB) (*UserStore, error) {
	store := &UserStore{
		db: db,
	}
	if err := store.db.Ping(); err != nil {
		return nil, err
	}
	return store, nil
}
