package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddContactTx(t *testing.T) {
	store := NewStore(testDB)

	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	fmt.Printf("user1: %#+v\n", user1)
	fmt.Printf("user2: %#+v\n", user2)

	var _type int16 = 0

	examine, err := testQueries.AddExamine(context.Background(), AddExamineParams{
		OwnerID:  user2.ID,
		TargetID: user1.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, examine)

	examines, err := testQueries.GetExamine(context.Background(), examine.OwnerID)
	require.NoError(t, err)
	require.NotEmpty(t, examines)

	for k := range examines {
		tc := examines[k]
		ct, err := store.AddContactTx(context.Background(), AddContactTxParams{
			AddContactParams: AddContactParams{
				OwnerID:  tc.OwnerID,
				TargetID: tc.TargetID,
				Type:     _type,
			},
			AfterCreate: func(contact Contact) error {
				c, err := store.AddContactTx(context.Background(), AddContactTxParams{
					AddContactParams: AddContactParams{
						OwnerID:  contact.TargetID,
						TargetID: contact.OwnerID,
						Type:     _type,
					},
					AfterCreate: func(contact Contact) error {
						fmt.Printf("contact: %#+v\n", contact)
						err := store.DeleteExamine(context.Background(), contact.TargetID)
						require.NoError(t, err)
						return err
					},
				})
				require.NoError(t, err)
				require.NotEmpty(t, c)
				return err
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, ct)
	}

}
