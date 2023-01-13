package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddContact(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	examine, err := testQueries.AddExamine(context.Background(), AddExamineParams{
		OwnerID:  user2.ID,
		TargetID: user1.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, examine)

	examines, err := testQueries.GetExamine(context.Background(), examine.OwnerID)
	require.NoError(t, err)
	require.NotEmpty(t, examines)

	for i := range examines {
		tc := examines[i]
		contact, err := testQueries.AddContact(context.Background(), AddContactParams{
			OwnerID:  tc.OwnerID,
			TargetID: tc.TargetID,
			Type:     0,
		})
		require.NoError(t, err)
		require.NotEmpty(t, contact)

		contact2, err := testQueries.AddContact(context.Background(), AddContactParams{
			OwnerID:  tc.TargetID,
			TargetID: tc.OwnerID,
			Type:     0,
		})
		require.NoError(t, err)
		require.NotEmpty(t, contact2)
	}

}
