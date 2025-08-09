package tests

import (
	"github.com/stretchr/testify/assert"
	"idm/inner/employee"
	"idm/inner/role"
	"testing"
)

func TestEmployeeRepository(t *testing.T) {
	a := assert.New(t)
	db := GetTestDataBase()
	clearDb := func() {
		db.MustExec("DELETE FROM employee")
	}

	defer func() {
		if r := recover(); r != nil {
			clearDb()
		}
	}()

	employeeRepository := employee.NewEmployeeRepository(db)
	fixture := NewFixtureEmployeeRepository(employeeRepository)

	t.Run("Add new Employee", func(t *testing.T) {
		entity := FixtureEmployeeEntity("New Employee")

		got, err := employeeRepository.AddEmployee(entity)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(entity.Name, got.Name)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)

		clearDb()
	})

	t.Run("Find employee by Id", func(t *testing.T) {
		employeeName := "NewEmployee"
		employeeId := fixture.Employee(employeeName)

		got, err := employeeRepository.FindById(employeeId)

		a.Nil(err)
		a.NotEmpty(got)
		a.NotEmpty(got.Id)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)
		a.Equal(got.Name, employeeName)

		clearDb()
	})

	t.Run("Find all employees", func(t *testing.T) {
		employeeName1 := "NewEmployee1"
		employeeId1 := fixture.Employee(employeeName1)
		employeeName2 := "NewEmployee2"
		employeeId2 := fixture.Employee(employeeName2)
		employeeName3 := "NewEmployee3"
		employeeId3 := fixture.Employee(employeeName3)
		_ = employeeRepository.DeleteByIdSilent(employeeId2)

		got, err := employeeRepository.FindAll()

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.Equal(employeeId1, got[0].Id)
		a.Equal(employeeName1, got[0].Name)
		a.Equal(employeeId3, got[1].Id)
		a.Equal(employeeName3, got[1].Name)

		clearDb()
	})

	t.Run("Find employers by ids", func(t *testing.T) {
		employeeName1 := "NewEmployee1"
		employeeId1 := fixture.Employee(employeeName1)
		employeeName2 := "NewEmployee2"
		_ = fixture.Employee(employeeName2)
		employeeName3 := "NewEmployee3"
		employeeId3 := fixture.Employee(employeeName3)
		employeeName4 := "NewEmployee4"
		employeeId4 := fixture.Employee(employeeName4)
		ids := []int64{employeeId1, employeeId3, employeeId4}
		_ = employeeRepository.DeleteByIdSilent(employeeId4)

		got, err := employeeRepository.FindByIds(ids)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.Equal(employeeId1, got[0].Id)
		a.Equal(employeeName1, got[0].Name)
		a.Equal(employeeId3, got[1].Id)
		a.Equal(employeeName3, got[1].Name)

		clearDb()
	})

	t.Run("Delete employee by Id", func(t *testing.T) {
		employeeName := "NewEmployee"
		employeeId := fixture.Employee(employeeName)

		err := employeeRepository.DeleteByIdSilent(employeeId)
		got, errFindById := employeeRepository.FindById(employeeId)

		a.Nil(err)
		a.ErrorIs(errFindById, employee.NotFound)
		a.Empty(got)

		clearDb()
	})

	t.Run("Delete all employees by Ids", func(t *testing.T) {
		employeeName1 := "NewEmployee1"
		employeeId1 := fixture.Employee(employeeName1)
		employeeName2 := "NewEmployee2"
		employeeId2 := fixture.Employee(employeeName2)
		ids := []int64{employeeId1, employeeId2}

		err := employeeRepository.DeleteByIdsSilent(ids)
		got, errFindByIds := employeeRepository.FindByIds(ids)

		a.Nil(err)
		a.ErrorIs(errFindByIds, employee.NotFound)
		a.Empty(got)

		clearDb()
	})

	t.Run("Delete employers by Ids", func(t *testing.T) {
		employeeName1 := "NewEmployee1"
		employeeId1 := fixture.Employee(employeeName1)
		employeeName2 := "NewEmployee2"
		employeeId2 := fixture.Employee(employeeName2)
		employeeName3 := "NewEmployee3"
		employeeId3 := fixture.Employee(employeeName3)
		ids := []int64{employeeId1, employeeId3}

		err := employeeRepository.DeleteByIdsSilent(ids)
		got, errFindByIds := employeeRepository.FindAll()

		a.Nil(err)
		a.Nil(errFindByIds)
		a.NotEmpty(got)
		a.Equal(1, len(got))
		a.Equal(employeeId2, got[0].Id)

		clearDb()
	})
}

func TestRoleRepository(t *testing.T) {
	a := assert.New(t)
	db := GetTestDataBase()
	clearDb := func() {
		db.MustExec("DELETE FROM role")
	}

	defer func() {
		if err := recover(); err != nil {
			clearDb()
		}
	}()

	roleRepository := role.NewRoleRepository(db)
	fixture := NewFixtureRoleRepository(roleRepository)

	t.Run("Add role", func(t *testing.T) {
		newRole := "New role"
		entity := FixtureRoleEntity(newRole)

		got, err := roleRepository.AddRole(entity)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(entity.Name, got.Name)
		a.NotEmpty(got.Id)
		a.True(got.Id > 0)
		a.NotEmpty(got.CreatedAt)
		a.NotEmpty(got.UpdatedAt)

		clearDb()
	})

	t.Run("Find role by Id", func(t *testing.T) {
		roleName := "New role"
		roleId := fixture.FixtureAddRole(roleName)

		got, err := roleRepository.FindById(roleId)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(roleName, got.Name)

		clearDb()
	})

	t.Run("Find all role", func(t *testing.T) {
		roleName1 := "New role 1"
		roleId1 := fixture.FixtureAddRole(roleName1)
		roleName2 := "New role 2"
		roleId2 := fixture.FixtureAddRole(roleName2)
		roleName3 := "New role 3"
		roleId3 := fixture.FixtureAddRole(roleName3)
		_ = roleRepository.DeleteByIdSilent(roleId2)

		got, err := roleRepository.FindAll()

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.Equal(roleId1, got[0].Id)
		a.Equal(roleName1, got[0].Name)
		a.Equal(roleId3, got[1].Id)
		a.Equal(roleName3, got[1].Name)

		clearDb()
	})

	t.Run("Find role by Ids", func(t *testing.T) {
		roleName1 := "New role 1"
		roleId1 := fixture.FixtureAddRole(roleName1)
		roleName2 := "New role 2"
		roleId2 := fixture.FixtureAddRole(roleName2)
		roleName3 := "New role 3"
		roleId3 := fixture.FixtureAddRole(roleName3)
		roleName4 := "New role 4"
		_ = fixture.FixtureAddRole(roleName4)
		_ = roleRepository.DeleteByIdSilent(roleId2)
		ids := []int64{roleId1, roleId2, roleId3}

		got, err := roleRepository.FindByIds(ids)

		a.Nil(err)
		a.NotEmpty(got)
		a.Equal(2, len(got))
		a.Equal(roleId1, got[0].Id)
		a.Equal(roleName1, got[0].Name)
		a.Equal(roleId3, got[1].Id)
		a.Equal(roleName3, got[1].Name)

		clearDb()
	})

	t.Run("Delete role by Id", func(t *testing.T) {
		roleName := "New role"
		roleId := fixture.FixtureAddRole(roleName)

		err := roleRepository.DeleteByIdSilent(roleId)
		got, errFindById := roleRepository.FindById(roleId)

		a.Nil(err)
		a.ErrorIs(errFindById, role.NotFound)
		a.Empty(got)

		clearDb()
	})

	t.Run("Delete role by Ids", func(t *testing.T) {
		roleName1 := "New role 1"
		roleId1 := fixture.FixtureAddRole(roleName1)
		roleName2 := "New role 2"
		roleId2 := fixture.FixtureAddRole(roleName2)
		roleName3 := "New role 3"
		roleId3 := fixture.FixtureAddRole(roleName3)
		ids := []int64{roleId1, roleId3}

		err := roleRepository.DeleteByIdsSilent(ids)
		got, errFindAll := roleRepository.FindAll()

		a.Nil(err)
		a.Nil(errFindAll)
		a.NotEmpty(got)
		a.Equal(1, len(got))
		a.Equal(roleId2, got[0].Id)

		clearDb()
	})
}
