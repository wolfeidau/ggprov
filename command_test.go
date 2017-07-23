package ggprov

import "testing"
import "github.com/stretchr/testify/assert"
import "os"

func TestRunCommandErrorStatus(t *testing.T) {
	err := RunCommand("id", []string{"-u", "bob"})
	assert.NotEmpty(t, err)
}

func TestRunCommand(t *testing.T) {
	user := os.Getenv("USER")
	err := RunCommand("id", []string{"-u", user})
	assert.Empty(t, err)
}

func TestUserNotExists(t *testing.T) {
	exists, err := UserExist("bob")
	assert.Empty(t, err)
	assert.False(t, exists)
}

func TestUserExists(t *testing.T) {
	user := os.Getenv("USER")
	exists, err := UserExist(user)
	assert.Empty(t, err)
	assert.True(t, exists)
}
