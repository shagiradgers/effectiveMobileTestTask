package store

import (
	"context"
	"database/sql"
	"effectiveMobileTestTask/internal/entity"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

const (
	CtxTimeout = 1 * time.Second
)

// UserParamsToFilter - Параметры для фильтрации пользователя
type UserParamsToFilter struct {
	Id          int
	Age         int
	Name        string
	Surname     string
	Patronymic  string
	Nationality string
	Sex         string
	Limit       uint
	Page        uint
}

type UserParamsToAdd struct {
	Name        interface{}
	Surname     interface{}
	Patronymic  interface{}
	Sex         interface{}
	Nationality interface{}
	Age         interface{}
}

type UserParamsToEdit struct {
	Name        sql.NullString
	Surname     sql.NullString
	Patronymic  sql.NullString
	Sex         sql.NullString
	Nationality sql.NullString
	Age         sql.NullInt16
}

type UserStore struct {
	db *sqlx.DB
}

func (u *UserStore) GetUsers(ctx context.Context, params UserParamsToFilter) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	var users []*entity.User
	sqlBuilder := squirrel.
		Select("id, name, surname, patronymic, age, sex, nationality").
		From("users")

	if params.Id != 0 {
		sqlBuilder = sqlBuilder.Where("id = ?", params.Id)
	}
	if len(strings.TrimSpace(params.Name)) != 0 {
		sqlBuilder = sqlBuilder.Where("name = ?", params.Name)
	}
	if len(strings.TrimSpace(params.Surname)) != 0 {
		sqlBuilder = sqlBuilder.Where("surname = ?", params.Surname)
	}
	if len(strings.TrimSpace(params.Patronymic)) != 0 {
		sqlBuilder = sqlBuilder.Where("patronymic = ?", params.Patronymic)
	}
	if len(strings.TrimSpace(params.Nationality)) != 0 {
		sqlBuilder = sqlBuilder.Where("nationality = ?", params.Nationality)
	}
	if len(strings.TrimSpace(params.Sex)) != 0 {
		sqlBuilder = sqlBuilder.Where("sex = ?", params.Sex)
	}
	if params.Age != 0 {
		sqlBuilder = sqlBuilder.Where("age = ?", params.Age)
	}
	if params.Limit != 0 {
		sqlBuilder = sqlBuilder.Limit(uint64(params.Limit))
	}
	if params.Page != 0 && params.Limit != 0 {
		sqlBuilder = sqlBuilder.Offset(uint64(params.Page * params.Limit))
	}
	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = u.db.SelectContext(ctx, users, query, args)
	if err != nil {
		return nil, err
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

	_, err = u.db.ExecContext(ctx, query, args)
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

	if params.Name.Valid {
		sqlBuilder = sqlBuilder.Set("name", params.Name.String)
	}
	if params.Surname.Valid {
		sqlBuilder = sqlBuilder.Set("surname", params.Surname.String)
	}
	if params.Patronymic.Valid {
		sqlBuilder = sqlBuilder.Set("patronymic", params.Patronymic.String)
	}
	if params.Nationality.Valid {
		sqlBuilder = sqlBuilder.Set("nationality", params.Nationality.String)
	}
	if params.Sex.Valid {
		sqlBuilder = sqlBuilder.Set("sex", params.Sex.String)
	}
	if params.Age.Valid {
		sqlBuilder = sqlBuilder.Set("age", params.Age.Int16)
	}

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) AddUser(ctx context.Context, name interface{}, surname interface{}, patronymic interface{}, sex interface{}, nationality interface{}, age interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	sqlBuilder := squirrel.
		Insert("users").
		Columns("name, surname, patronymic, nationality, sex, age")
	//Values(name.(string), surname.(string), patronymic.(string), nationality.(string), sex.(string), age.(int))

	//TODO: доработать
	if v, ok := name.(string); ok {
		sqlBuilder = sqlBuilder.Values(v)
	} else {
		sqlBuilder = sqlBuilder.Values(sql.NullString{})
	}

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.ExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func New(db *sqlx.DB) (*UserStore, error) {
	store := &UserStore{
		db: db,
	}
	if err := store.db.Ping(); err != nil {
		return nil, err
	}
	return store, nil
}
