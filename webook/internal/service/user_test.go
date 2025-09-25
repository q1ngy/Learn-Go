package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("abc123##")
	fromPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	err = bcrypt.CompareHashAndPassword(fromPassword, []byte("abc123##"))
	assert.NoError(t, err)
}
