package utils

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

	t.Run("No Error", func(t *testing.T) {
		err := ValidateUsername("testmailcom_")
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := ValidateUsername("testmailcom_1111111111")
		require.NotEmpty(t, err)
	})
}
