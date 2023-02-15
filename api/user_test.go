package api

import (
	mockdb "IMChat/db/pg/mock"
	db "IMChat/db/pg/sqlc"
	"IMChat/token"
	"IMChat/utils"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// func setCache(t *testing.T, cache cache.Cache, user db.User) {
// 	key := fmt.Sprintf("user:%d_%s", user.ID, user.Username)
// 	err := cache.SetCache(context.Background(), key, user)
// 	require.NoError(t, err)
// }

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"email":    user.Email,
				"nickname": user.Nickname,
				"gender":   user.Gender,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					Password: password,
					Email:    user.Email,
					Nickname: user.Nickname,
					Gender:   user.Gender,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"email":    user.Email,
				"nickname": user.Nickname,
				"gender":   user.Gender,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-userename",
				"password": password,
				"email":    user.Email,
				"nickname": user.Nickname,
				"gender":   user.Gender,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"email":    "invalid-email",
				"nickname": user.Nickname,
				"gender":   user.Gender,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "InvalidGender",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"email":    user.Email,
				"nickname": user.Nickname,
				"gender":   0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"email":    user.Email,
				"nickname": user.Nickname,
				"gender":   3,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recoder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username": user.Username,
				"password": "123",
				"email":    user.Email,
				"nickname": user.Nickname,
				"gender":   3,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			cache := newTestRedis(t)

			server := newTestServer(t, store, cache)
			recoder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(recoder)
		})
	}

}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "UsernameAndPasswordOK",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				require.NotEmpty(t, recoder.Body)
			},
		},
		{
			name: "UsernameAndEmailOK",
			body: gin.H{
				"username": user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				require.NotEmpty(t, recoder.Body)
			},
		},
		{
			name: "UserNotFund",
			body: gin.H{
				"username": "NotFound",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code)
			},
		},
		{
			name: "IncorrectPassword",
			body: gin.H{
				"username": user.Username,
				"password": "incorrect",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
			},
		},
		{
			name: "BadRequestError",
			body: gin.H{
				"username": "",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			cache := newTestRedis(t)

			server := newTestServer(t, store, cache)
			recoder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.checkResponse(recoder)
		})
	}
}

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
}

func (e eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, arg.Password.String)
	if err != nil {
		return false
	}

	e.arg.Password.String = arg.Password.String
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateUserParams(arg db.UpdateUserParams, password string) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg, password}
}

func TestUpdateUser(t *testing.T) {
	user, password := randomUser(t)

	email := utils.RandomEmail()
	nickname := utils.RandomString(6)
	var gender int16 = 1
	avatar := fmt.Sprintf("%s.jpg", utils.RandomString(6))

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.Maker)
		buildStubd    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"email":    utils.RandomEmail(),
				"password": password,
				"nickname": utils.RandomString(6),
				"gender":   1,
				"avatar":   fmt.Sprintf("%s.jpg", utils.RandomString(6)),
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorizatin(t, req, tokenMaker, user.Username, time.Minute)
			},
			buildStubd: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UpdateUserAllFields",
			body: gin.H{
				"username": user.Username,
				"email":    email,
				"password": password,
				"nickname": nickname,
				"gender":   gender,
				"avatar":   avatar,
			},
			setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.Maker) {
				addAuthorizatin(t, req, tokenMaker, user.Username, time.Minute)
			},
			buildStubd: func(store *mockdb.MockStore) {
				hashPassword, err := utils.HashPassword(password)
				require.NoError(t, err)

				arg := db.UpdateUserParams{
					Username: user.Username,
					Password: sql.NullString{
						String: hashPassword,
						Valid:  true,
					},
					Email: sql.NullString{
						String: email,
						Valid:  true,
					},
					Nickname: sql.NullString{
						String: nickname,
						Valid:  true,
					},
					Gender: sql.NullInt16{
						Int16: gender,
						Valid: true,
					},
					Avatar: sql.NullString{
						String: avatar,
						Valid:  true,
					},
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), EqUpdateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubd(store)

			cache := newTestRedis(t)

			server := newTestServer(t, store, cache)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(6)

	hashPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	str := utils.RandomString(6)

	user = db.User{
		Username: str,
		Password: hashPassword,
		Email:    utils.RandomEmail(),
		Nickname: str,
		Gender:   3,
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Gender, gotUser.Gender)
	require.Empty(t, gotUser.Password)
}
