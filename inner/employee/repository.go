package employee

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
	NotFound = errors.New("employee not found")
)

func NewEmployeeRepository(database *sqlx.DB) *Repository {
	return &Repository{db: database}
}

type Entity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// AddEmployee - добавить новый элемент в коллекцию
func (r *Repository) AddEmployee(employee *Entity) (e Entity, err error) {
	q := "INSERT INTO employee(name) VALUES($1) RETURNING id, created_at, updated_at, name"
	err = r.db.Get(&e, q, employee.Name)
	return e, err
}

// FindById - найти элемент коллекции по его id (этот метод мы реализовали на уроке)
func (r *Repository) FindById(id int64) (e Entity, err error) {
	q := "SELECT id, name, created_at, updated_at FROM employee e WHERE e.is_deleted = FALSE AND e.id=$1"
	if err = r.db.Get(&e, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, NotFound
		}
		return e, fmt.Errorf("employeeRepository.FindById: %w", err)
	}
	return e, err
}

// FindAll - найти все элементы коллекции
func (r *Repository) FindAll() (employees []Entity, err error) {
	q := "SELECT  id, name, created_at, updated_at  FROM employee e WHERE e.is_deleted = FALSE"
	if err = r.db.Select(&employees, q); err != nil {
		return employees, fmt.Errorf("employeeRepository.FindAll: %w", err)
	}
	if len(employees) == 0 {
		return employees, NotFound
	}
	return employees, err
}

// FindByIds - найти слайс элементов коллекции по слайсу их id
func (r *Repository) FindByIds(ids []int64) (employees []Entity, err error) {
	q := "SELECT  id, name, created_at, updated_at  FROM employee e WHERE e.is_deleted = FALSE AND e.id IN (?)"
	query, args, errQueryBuild := sqlx.In(q, ids)
	if errQueryBuild != nil {
		return nil, errQueryBuild
	}
	query = sqlx.Rebind(2, query)
	if err = r.db.Select(&employees, query, args...); err != nil {
		return employees, fmt.Errorf("employeeRepository.FindByIds: %w", err)
	}
	if len(employees) == 0 {
		return employees, NotFound
	}
	return employees, err
}

// DeleteByIdSilent - удалить элемент коллекции по его id
func (r *Repository) DeleteByIdSilent(id int64) (err error) {
	q := "UPDATE employee e SET is_deleted = TRUE WHERE e.is_deleted = FALSE and e.id = $1"
	_, err = r.db.Exec(q, id)
	return err
}

// DeleteByIdsSilent - удалить элементы по слайсу их id
func (r *Repository) DeleteByIdsSilent(ids []int64) (err error) {
	q := "UPDATE employee e SET is_deleted = TRUE WHERE e.is_deleted = FALSE and e.id IN (?)"
	query, args, errQueryBuild := sqlx.In(q, ids)
	if errQueryBuild != nil {
		return errQueryBuild
	}
	query = sqlx.Rebind(2, query)
	_, err = r.db.Exec(query, args...)
	return err
}
