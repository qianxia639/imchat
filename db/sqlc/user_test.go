package db

import (
	"IMChat/utils"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	name := utils.RandomString(6)
	pwd := utils.RandomString(6)
	hashPwd, err := utils.HashPassword(pwd)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: name,
		Email:    utils.RandomEmail(),
		Nickname: name,
		Password: hashPwd,
		Gender:   3,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Nickname, user.Nickname)
	require.Equal(t, arg.Password, user.Password)
	// require.Equal(t, arg.Gender, user.Gender)
	// require.NotZero(t, user.Gender)
	// require.NotZero(t, user.Avatar)
	require.False(t, user.RegisterTime.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	testCase := []struct {
		name string
		Func func(*Queries) (User, error)
	}{
		{
			name: "Get User Username",
			Func: func(q *Queries) (User, error) {
				return q.GetUser(context.Background(), user1.Username)
			},
		},
		{
			name: "Get User Email",
			Func: func(q *Queries) (User, error) {
				return q.GetUser(context.Background(), user1.Username)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			user2, err := tc.Func(testQueries)
			require.NoError(t, err)
			require.NotEmpty(t, user2)

			require.Equal(t, user1.ID, user2.ID)
			require.Equal(t, user1.Username, user2.Username)
			require.Equal(t, user1.Email, user2.Email)
			require.Equal(t, user1.Nickname, user2.Nickname)
			require.Equal(t, user1.Password, user2.Password)
			// require.Equal(t, user1.Gender, user2.Gender)
			// require.Equal(t, user1.Avatar, user2.Avatar)

			require.WithinDuration(t, user1.RegisterTime, user2.RegisterTime, time.Second)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user1 := createRandomUser(t)

	testCase := []struct {
		name string
		fn   func(*Queries) (User, error)
	}{
		{
			name: "AND Password",
			fn: func(q *Queries) (User, error) {
				return q.LoginUser(context.Background(), LoginUserParams{
					Username: user1.Username,
					Password: user1.Password,
				})
			},
		},
		{
			name: "Email AND Password",
			fn: func(q *Queries) (User, error) {
				return q.LoginUser(context.Background(), LoginUserParams{
					Username: user1.Email,
					Password: user1.Password,
				})
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			user2, err := tc.fn(testQueries)
			require.NoError(t, err)
			require.NotEmpty(t, user2)

			require.Equal(t, user1.ID, user2.ID)
			require.Equal(t, user1.Username, user2.Username)
			require.Equal(t, user1.Email, user2.Email)
			require.Equal(t, user1.Nickname, user2.Nickname)
			require.Equal(t, user1.Password, user2.Password)
			// require.Equal(t, user1.Gender, user2.Gender)
			// require.Equal(t, user1.Avatar, user2.Avatar)

			require.WithinDuration(t, user1.RegisterTime, user2.RegisterTime, time.Second)
		})
	}
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := utils.RandomEmail()
	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEqual(t, oldUser.Email, user.Email)
	require.Equal(t, newEmail, user.Email)
	require.Equal(t, oldUser.Username, user.Username)
}

func TestUpdateUserOnlyNickname(t *testing.T) {
	oldUser := createRandomUser(t)

	newNickname := utils.RandomEmail()
	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Nickname: sql.NullString{
			String: newNickname,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEqual(t, oldUser.Nickname, user.Nickname)
	require.Equal(t, newNickname, user.Nickname)
	require.Equal(t, oldUser.Username, user.Username)
	require.Equal(t, oldUser.Password, user.Password)
	require.Equal(t, oldUser.Gender, user.Gender)
	require.Equal(t, oldUser.Avatar, user.Avatar)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := utils.RandomString(6)
	newHashPassword, err := utils.HashPassword(newPassword)
	require.NoError(t, err)

	_, err = testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Password: sql.NullString{
			String: newHashPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
}

func TestUpdateUserOnlyGender(t *testing.T) {
	oldUser := createRandomUser(t)

	newGender := 1
	_, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Gender: sql.NullInt16{
			Int16: int16(newGender),
			Valid: true,
		},
	})
	require.NoError(t, err)
}

func TestUpdateUserOnlyAvatar(t *testing.T) {
	oldUser := createRandomUser(t)

	newAvatar := "https://ns-strategy.cdn.bcebos.com/ns-strategy/upload/fc_big_pic/part-00544-3205.jpg"
	_, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Avatar: sql.NullString{
			String: newAvatar,
			Valid:  true,
		},
	})
	require.NoError(t, err)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := utils.RandomEmail()
	newNickname := utils.RandomString(6)
	var newGender int16 = 1
	newPassword := utils.RandomString(6)
	newHashPassword, err := utils.HashPassword(newPassword)
	require.NoError(t, err)

	newAvatar := "https://ns-strategy.cdn.bcebos.com/ns-strategy/upload/fc_big_pic/part-00544-3215.jpg"
	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		Nickname: sql.NullString{
			String: newNickname,
			Valid:  true,
		},
		Password: sql.NullString{
			String: newHashPassword,
			Valid:  true,
		},
		Gender: sql.NullInt16{
			Int16: int16(newGender),
			Valid: true,
		},
		Avatar: sql.NullString{
			String: newAvatar,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	require.Equal(t, newEmail, user.Email)
	require.Equal(t, newNickname, user.Nickname)
	require.Equal(t, newHashPassword, user.Password)
	require.Equal(t, newGender, user.Gender)
	require.Equal(t, newAvatar, user.Avatar)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
}
