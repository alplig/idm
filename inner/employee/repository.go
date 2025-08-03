package employee

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(database *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{db: database}
}

type EmployeeEntity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
}

// AddEmployee - добавить новый элемент в коллекцию
func (r *EmployeeRepository) AddEmployee(name string) (employee *EmployeeEntity, err error) {
	err = r.db.Get(&employee, "INSERT INTO employee(name) VALUES(?) RETURNING *", name)
	return employee, err
}

// FindById - найти элемент коллекции по его id (этот метод мы реализовали на уроке)
func (r *EmployeeRepository) FindById(id int64) (employee *EmployeeEntity, err error) {
	err = r.db.Get(&employee, "SELECT * FROM employee e WHERE e.is_deleted = FALSE AND e.id=?", id)
	return employee, err
}

// FindAll - найти все элементы коллекции
func (r *EmployeeRepository) FindAll() (employees []*EmployeeEntity, err error) {
	err = r.db.Select(&employees, "SELECT * FROM employee e WHERE e.is_deleted = FALSE")
	return employees, err
}

// FindByIds - найти слайс элементов коллекции по слайсу их id
func (r *EmployeeRepository) FindByIds(ids []int64) (employees []*EmployeeEntity, err error) {
	query, args, errQueryBuild := sqlx.In(
		"SELECT * FROM employee e WHERE e.is_deleted = FALSE AND e.id IN (?)",
		ids)
	if errQueryBuild != nil {
		return nil, errQueryBuild
	}
	err = r.db.Select(&employees, query, args...)
	return employees, err
}

// DeleteByIdSilent - удалить элемент коллекции по его id
func (r *EmployeeRepository) DeleteByIdSilent(id int64) (err error) {
	_, err = r.db.Exec("UPDATE employee e SET is_deleted = TRUE WHERE e.is_deleted = FALSE and e.id = ?", id)
	return err
}

// DeleteByIdsSilent - удалить элементы по слайсу их id
func (r *EmployeeRepository) DeleteByIdsSilent(ids []int64) (err error) {
	query, args, errQueryBuild := sqlx.In(
		"UPDATE employee e SET is_deleted = TRUE WHERE e.is_deleted = FALSE and e.id IN (?)",
		ids)
	if errQueryBuild != nil {
		return errQueryBuild
	}
	_, err = r.db.Exec(query, args...)
	return err
}
