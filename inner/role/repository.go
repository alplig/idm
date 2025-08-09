package role

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Repository struct {
	db *sqlx.DB
}

var (
	NotFound = errors.New("role not found")
)

func NewRoleRepository(database *sqlx.DB) *Repository {
	return &Repository{db: database}
}

type Entity struct {
	Id        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
}

// AddRole - добавить новый элемент в коллекцию
func (rr *Repository) AddRole(role *Entity) (r Entity, err error) {
	q := "INSERT INTO role(name) VALUES($1) RETURNING id, created_at, updated_at, name"
	err = rr.db.Get(&r, q, role.Name)
	return r, err
}

// FindById - найти элемент коллекции по его id (этот метод мы реализовали на уроке)
func (rr *Repository) FindById(id int64) (role Entity, err error) {
	q := "SELECT id, created_at, updated_at, name FROM role rl WHERE rl.is_deleted = FALSE AND rl.id=$1"
	if err = rr.db.Get(&role, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entity{}, NotFound
		}
		return Entity{}, fmt.Errorf("role.FindById: %w", err)
	}
	return role, err
}

// FindAll - найти все элементы коллекции
func (rr *Repository) FindAll() (roles []Entity, err error) {
	q := "SELECT id, created_at, updated_at, name FROM role rl WHERE rl.is_deleted = FALSE"
	err = rr.db.Select(&roles, q)
	return roles, err
}

// FindByIds - найти слайс элементов коллекции по слайсу их id
func (rr *Repository) FindByIds(ids []int64) (roles []Entity, err error) {
	q := "SELECT id, created_at, updated_at, name FROM role rl WHERE rl.is_deleted = FALSE AND rl.id IN (?)"
	query, args, errQueryBuild := sqlx.In(q, ids)
	query = sqlx.Rebind(2, query)
	if errQueryBuild != nil {
		return nil, errQueryBuild
	}
	err = rr.db.Select(&roles, query, args...)
	return roles, err
}

// DeleteByIdSilent - удалить элемент коллекции по его id
func (rr *Repository) DeleteByIdSilent(id int64) (err error) {
	q := "UPDATE role rl SET is_deleted = TRUE WHERE rl.is_deleted = FALSE and rl.id = $1"
	_, err = rr.db.Exec(q, id)
	return err
}

// DeleteByIdsSilent - удалить элементы по слайсу их id
func (rr *Repository) DeleteByIdsSilent(ids []int64) (err error) {
	q := "UPDATE role rl SET is_deleted = TRUE WHERE rl.is_deleted = FALSE and rl.id IN (?)"
	query, args, errQueryBuild := sqlx.In(q, ids)
	query = sqlx.Rebind(2, query)
	if errQueryBuild != nil {
		return errQueryBuild
	}
	_, err = rr.db.Exec(query, args...)
	return err
}
