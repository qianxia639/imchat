package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddExamine(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	examine, err := testQueries.AddExamine(context.Background(), AddExamineParams{
		OwnerID:  user1.ID,
		TargetID: user2.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, examine)

}

func TestGetExamine(t *testing.T) {

	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	examine, err := testQueries.AddExamine(context.Background(), AddExamineParams{
		OwnerID:  user1.ID,
		TargetID: user2.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, examine)

	examines, err := testQueries.GetExamine(context.Background(), examine.OwnerID)
	require.NoError(t, err)
	require.NotEmpty(t, examines)
}
