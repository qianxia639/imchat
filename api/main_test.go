package api

import (
	db "IMChat/db/pg/sqlc"
	"IMChat/utils"
	"IMChat/utils/config"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	conf := config.Config{}
	conf.Token.TokenSymmetricKey = utils.RandomString(32)
	conf.Token.AccessTokenDuration = time.Minute

	server, err := NewServer(conf, store)
	require.NoError(t, err)

	return server
}

func newTestDb(t *testing.T) db.Store {
	conf, err := config.LoadConfig("../.")
	require.NoError(t, err)

	conn, err := sql.Open(conf.Postgres.Driver, conf.Postgres.Source)
	require.NoError(t, err)

	return db.NewStore(conn)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
