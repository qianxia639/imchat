package db

import "context"

type AddContactTxParams struct {
	AddContactParams
	AfterCreate func(contact Contact) error
}

type AddContactxResult struct {
	Contact Contact
}

func (store *SQLStore) AddContactTx(ctx context.Context, arg AddContactTxParams) (AddContactxResult, error) {
	var result AddContactxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Contact, err = q.AddContact(ctx, arg.AddContactParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Contact)
	})

	return result, err
}
