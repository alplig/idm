package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"idm/inner/common"
	"idm/inner/database"
	"os"
	"path/filepath"
	"testing"
)

// 1  в проекте нет .env  файла (должны получить конфигурацию из пременных окружения)
func Test_getConfigEnvFileNotExistEnvVarExist(t *testing.T) {
	a := assert.New(t)
	expDriverName := "moc_driver_name"
	expDSN := "moc_DSN"
	t.Setenv("DB_DRIVER_NAME", expDriverName)
	t.Setenv("DB_DSN", expDSN)
	expected := common.Config{DbDriverName: expDriverName, Dsn: expDSN}
	actual := common.GetConfig("")
	a.EqualExportedValues(expected, actual)
}

// 2 в проекте есть .env  файл, но в нём нет нужных переменных и в переменных окружения их тоже нет
// (должны получить пустую структуру idm.inner.common.Config)
func Test_getConfigEnvFileNotExistEnvVarNotExist(t *testing.T) {
	a := assert.New(t)
	actual := common.GetConfig("")
	a.Empty(actual)
}

// 3 в проекте нет .env  файла (должны получить конфигурацию из пременных окружения)
func Test_getConfigEnvFileEmptyEnvVarExist(t *testing.T) {
	a := assert.New(t)
	tmpDir := t.TempDir()
	envTestPath := filepath.Join(tmpDir, ".env_test")
	if envFile, err := os.Create(envTestPath); err != nil {
		t.Fatalf("Не удалось создать тестовый файл .env_test: %s", err)
	} else {
		envFile.Close()
	}
	expDriverName := "moc_driver_name"
	expDSN := "moc_DSN"
	t.Setenv("DB_DRIVER_NAME", expDriverName)
	t.Setenv("DB_DSN", expDSN)
	expected := common.Config{DbDriverName: expDriverName, Dsn: expDSN}
	actual := common.GetConfig(envTestPath)
	a.EqualExportedValues(expected, actual)
}

// 4 в проекте есть корректно заполненный .env файл, в переменных окружения нет конфликтующих с ним переменных
// (должны получить структуру  idm.inner.common.Config, заполненную данными из .env файла)
func Test_getConfigEnvFileActualEnvVarEmpty(t *testing.T) {
	a := assert.New(t)
	expDriverName := "moc_driver_name"
	expDSN := "moc_DSN"
	tmpDir := t.TempDir()
	envTestPath := filepath.Join(tmpDir, ".env_test")
	envFileText := fmt.Sprintf("DB_DRIVER_NAME=%s\nDB_DSN=%s", expDriverName, expDSN)
	if err := os.WriteFile(envTestPath, []byte(envFileText), 0o600); err != nil {
		t.Fatalf("Не удалось создать тестовый файл .env_test: %s", err)
	}
	expected := common.Config{DbDriverName: expDriverName, Dsn: expDSN}
	actual := common.GetConfig(envTestPath)
	a.EqualExportedValues(expected, actual)
}

// 5 в проекте есть .env  файл и в нём есть нужные переменные, но в переменных окружения они тоже есть
// (с другими значениями) - должны получить структуру  idm.inner.common.Config, заполненную данными.
// Нужно проверить, какими значениями она будет заполнена (из .env файла или из переменных окружения)
func Test_getConfigEnvFileExistEnvVarExist(t *testing.T) {
	a := assert.New(t)
	expDriverName := "moc_driver_name"
	expDSN := "moc_DSN"
	tmpDir := t.TempDir()
	envTestPath := filepath.Join(tmpDir, ".env_test")
	envFileText := "DB_DRIVER_NAME=file_driver_name\nDB_DSN=file_dsn"
	if err := os.WriteFile(envTestPath, []byte(envFileText), 0o600); err != nil {
		t.Fatalf("Не удалось создать тестовый файл .env_test: %s", err)
	}
	t.Setenv("DB_DRIVER_NAME", expDriverName)
	t.Setenv("DB_DSN", expDSN)
	expected := common.Config{DbDriverName: expDriverName, Dsn: expDSN}
	actual := common.GetConfig(envTestPath)
	a.EqualExportedValues(expected, actual)
}

// 6 приложение не может подключиться к базе данных с некорректным конфигом
// (например, неправильно указан: хост, порт, имя базы данных, логин или пароль)
func Test_connectDbWrongCredentials(t *testing.T) {
	a := assert.New(t)
	cnf := common.Config{DbDriverName: "postgres", Dsn: "host=127.0.0.1 port=5432 user=admin password=4323 dbname=idm sslmode=disable"}
	a.Panics(func() { database.ConnectDbWithCfg(cnf) })
}

// 7 приложение может подключиться к базе данных с корректным конфигом.
func Test_connectDbWithCfg(t *testing.T) {
	a := assert.New(t)
	c := database.ConnectDb()
	status := true
	if err := c.Ping(); err != nil {
		status = false
	}
	a.Truef(status, "Не удалось подключиться к базе данных")
}
