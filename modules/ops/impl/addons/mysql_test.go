package addons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	url      = "jdbc:mysql://127.0.0.1:3306"
	user     = "root"
	password = "password"
)

func TestParseAddr(t *testing.T) {
	host, port := ParseAddr(url)
	assert.Equal(t, host, "127.0.0.1")
	assert.Equal(t, port, "3306")
}

func TestGetDSN(t *testing.T) {
	dsn := GetDSN(url, user, password)
	assert.Equal(t, dsn, "root:password@tcp(127.0.0.1:3306)/?charset=utf8mb4,utf8&parseTime=true")
}
