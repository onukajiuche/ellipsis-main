// build+ unit
package utility_test

import (
	"brief/utility"
	"testing"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	hashPassword(t, "password")
}

func TestPasswordIsValid(t *testing.T) {

	password := "password"
	hash, salt := hashPassword(t, password)

	if ok := utility.PasswordIsValid(password, salt, hash); !ok {
		t.Errorf("Expected '%v', but got '%v'", !ok, ok)
	}
}

func TestGetURLHash(t *testing.T) {
	tests := []struct {
		Name        string
		ID          string
		URL         string
		ErrIsNil    bool
		HashIsEmpty bool
	}{
		{"Valid_URL_and_ID", uuid.NewString(), "https://randomUrl.com", true, false},
		{"Empty_URL_and_ID", "", "", true, false},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			hashV, err := utility.GetURLHash(test.ID, test.URL)
			if (err == nil) != test.ErrIsNil {
				t.Errorf("Expected 'error' not to be '%v'", err)
			}

			if (hashV == "") != test.HashIsEmpty {
				t.Errorf("Expected 'hash' not to be '%v'", hashV)
			}

		})
	}
}

func hashPassword(t *testing.T, password string) (hash, salt string) {
	var err error
	hash, salt, err = utility.HashPassword("password")
	if err != nil {
		t.Errorf("Expected 'error' to be nil, got 'error': %v", err)
	}

	if hash == "" {
		t.Errorf("Expected 'hash' not to be empty")
	}

	if salt == "" {
		t.Errorf("Expected 'salt' not to be empty")
	}

	return
}
