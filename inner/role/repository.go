package role

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(database *sqlx.DB) *RoleRepository {
	return &RoleRepository{db: database}
}

type RoleEntity struct {
	Id        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsDeleted bool      `db:"is_deleted"`
	Name      string    `db:"name"`
}

// AddRole - добавить новый элемент в коллекцию
func (rr *RoleRepository) AddRole(role *RoleEntity) (r RoleEntity, err error) {
	err = rr.db.Get(&r, "INSERT INTO role(name) VALUES(?) RETURNING *", role.Name)
	return r, err
}

// FindById - найти элемент коллекции по его id (этот метод мы реализовали на уроке)
func (rr *RoleRepository) FindById(id int64) (role RoleEntity, err error) {
	err = rr.db.Get(&role, "SELECT * FROM role rl WHERE rl.is_deleted = FALSE AND rl.id=?", id)
	return role, err
}

// FindAll - найти все элементы коллекции
func (rr *RoleRepository) FindAll() (roles []RoleEntity, err error) {
	err = rr.db.Select(&roles, "SELECT * FROM role rl WHERE rl.is_deleted = FALSE")
	return roles, err
}

// FindByIds - найти слайс элементов коллекции по слайсу их id
func (rr *RoleRepository) FindByIds(ids []int64) (roles []RoleEntity, err error) {
	query, args, errQueryBuild := sqlx.In("SELECT * FROM role rl WHERE rl.is_deleted = FALSE AND rl.id IN (?)", ids)
	if errQueryBuild != nil {
		return nil, errQueryBuild
	}
	err = rr.db.Select(&roles, query, args...)
	return roles, err
}

// DeleteByIdSilent - удалить элемент коллекции по его id
func (rr *RoleRepository) DeleteByIdSilent(id int64) (err error) {
	_, err = rr.db.Exec("UPDATE role rl SET is_deleted = TRUE WHERE rl.is_deleted = FALSE and rl.id = ?", id)
	return err
}

// DeleteByIdsSilent - удалить элементы по слайсу их id
func (rr *RoleRepository) DeleteByIdsSilent(ids []int64) (err error) {
	query, args, errQueryBuild := sqlx.In("UPDATE role rl SET is_deleted = TRUE WHERE rl.is_deleted = FALSE and rl.id IN (?)", ids)
	if errQueryBuild != nil {
		return errQueryBuild
	}
	_, err = rr.db.Exec(query, args...)
	return err
}
