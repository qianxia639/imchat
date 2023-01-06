package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	err := ValidateEmail("test@mail.com")
	require.NoError(t, err)
}

func TestValidateLen(t *testing.T) {
	err := ValidateLen("test@mail.com", 5, 20)
	require.NoError(t, err)
}

func TestValidateUsername(t *testing.T) {
	err := ValidateUsername("testmailcom_")
	require.NoError(t, err)
}
func TestValidateGender(t *testing.T) {
	err := ValidateGender(1)
	require.NoError(t, err)
}
