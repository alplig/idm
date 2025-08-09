package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/role"
	"idm/utils"
	"path/filepath"
)

type FixtureEmployee struct {
	employees *employee.Repository
}

func NewFixtureEmployeeRepository(employees *employee.Repository) *FixtureEmployee {
	return &FixtureEmployee{employees}
}

func FixtureEmployeeEntity(name string) *employee.Entity {
	return &employee.Entity{
		Name: name,
	}
}

func (fe *FixtureEmployee) Employee(name string) int64 {
	entity := FixtureEmployeeEntity(name)
	emp, err := fe.employees.AddEmployee(entity)
	if err != nil {
		panic(err)
	}
	return emp.Id
}

type FixtureRole struct {
	role *role.Repository
}

func NewFixtureRoleRepository(role *role.Repository) *FixtureRole {
	return &FixtureRole{role}
}

func (fr *FixtureRole) FixtureAddRole(name string) int64 {
	entity := FixtureRoleEntity(name)
	rl, err := fr.role.AddRole(entity)
	if err != nil {
		panic(err)
	}
	return rl.Id
}

func FixtureRoleEntity(name string) *role.Entity {
	return &role.Entity{
		Name: name,
	}
}

func GetTestDataBase() *sqlx.DB {
	db := database.ConnectDb(database.TestEnv)
	pathRoot, _ := utils.FindRoot()
	pathMigrations := filepath.Join(pathRoot, "./migrations")
	if err := goose.Up(db.DB, pathMigrations); err != nil {
		panic(err)
	}
	return db
}
