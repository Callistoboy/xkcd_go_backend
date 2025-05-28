package aaa

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoginFail(t *testing.T) {
	// Defining the columns of the table
	a := AAA{users: map[string]string{"admin": "password"}}
	var tests = []struct {
		name     string
		password string
	}{
		// the table itself
		{"test", ""},
		{"", ""},
		{"admin", "admin"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := a.Login(tt.name, tt.password, "")
			if ans != "" || err == nil {
				t.Errorf(`Login(%q, %q) = %v, %v, want "", error`, tt.name, tt.password, ans, err)
			}
		})
	}
}

func TestLoginOK(t *testing.T) {
	a := AAA{users: map[string]string{"admin": "password"}}
	_, err := a.Login("admin", "password", "")
	require.NoError(t, err)
}

type DummyLogger struct{}

func (d DummyLogger) Error(msg string, keysAndValues ...interface{}) {}

func TestVerifyOK(t *testing.T) {
	a := AAA{
		users:    map[string]string{"admin": "password"},
		log:      &DummyLogger{},
		tokenTTL: time.Hour,
	}
	token, err := a.Login("admin", "password", "")
	require.NoError(t, err)
	t.Logf("token: %s", token)
	err = a.Verify(token)
	require.NoError(t, err)
}

func TestVerifyNoTTL(t *testing.T) {
	a := AAA{
		users: map[string]string{"admin": "password"},
		log:   &DummyLogger{},
	}

	token, err := a.Login("admin", "password", "")
	require.NoError(t, err)
	err = a.Verify(token)
	require.Error(t, err)
}

func TestVerifyNoAdmin(t *testing.T) {
	a := AAA{
		users:    map[string]string{"admin": "password"},
		log:      &DummyLogger{},
		tokenTTL: time.Hour,
	}

	token, err := a.Login("admin", "password", "user")
	require.NoError(t, err)
	err = a.Verify(token)
	require.Error(t, err)
}

func TestVerifyFail(t *testing.T) {
	a := AAA{
		users:    map[string]string{"admin": "password"},
		log:      &DummyLogger{},
		tokenTTL: time.Hour,
	}

	var tests = []struct {
		token string
	}{
		{"test"},
		{"Token "},
		{"Token test"},
		{""},
	}
	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			err := a.Verify(tt.token)
			require.Error(t, err)
		})
	}
}

func TestNewOK(t *testing.T) {
	os.Setenv("ADMIN_USER", "admin")
	os.Setenv("ADMIN_PASSWORD", "password")

	_, err := New(time.Hour, &DummyLogger{})
	require.NoError(t, err)
	os.Clearenv()
}

func TestNewFail(t *testing.T) {
	_, err := New(time.Hour, &DummyLogger{})
	require.Error(t, err)

	os.Setenv("ADMIN_USER", "admin")
	_, err = New(time.Hour, &DummyLogger{})
	require.Error(t, err)

	os.Clearenv()
	os.Setenv("ADMIN_PASSWORD", "password")
	_, err = New(time.Hour, &DummyLogger{})
	require.Error(t, err)
	os.Clearenv()
}
