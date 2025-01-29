package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/tempura/pkg/contants/envvarname"
	"github.com/tecnologer/tempura/pkg/dao/db"
)

func TestNewConfigFromEnvVars(t *testing.T) {
	t.Run("no_env_var_data", func(t *testing.T) {
		t.Setenv(envvarname.DBHost, "")
		t.Setenv(envvarname.DBPort, "")
		t.Setenv(envvarname.DBUsername, "")
		t.Setenv(envvarname.DBPassword, "")
		t.Setenv(envvarname.DBName, "")
		t.Setenv(envvarname.DBSSLMode, "")

		want := &db.Config{
			Host:     db.DefaultHost,
			Port:     db.DefaultPort,
			User:     db.DefaultUser,
			Password: "",
			DBName:   db.DefaultDB,
			SSLMode:  db.DefaultSSLMode,
		}

		got := db.NewConfigFromEnvVars()
		assert.Equal(t, want, got)
	})

	t.Run("env_var_data", func(t *testing.T) {
		t.Setenv(envvarname.DBHost, "dbhost")
		t.Setenv(envvarname.DBPort, "1234")
		t.Setenv(envvarname.DBUsername, "user")
		t.Setenv(envvarname.DBPassword, "test_password")
		t.Setenv(envvarname.DBName, "test_db")
		t.Setenv(envvarname.DBSSLMode, "undefined")

		want := &db.Config{
			Host:     "dbhost",
			Port:     "1234",
			User:     "user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "undefined",
		}

		got := db.NewConfigFromEnvVars()
		assert.Equal(t, want, got)
	})
}
