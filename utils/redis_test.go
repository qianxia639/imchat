package utils

import (
	"IMChat/utils/config"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitRedis(t *testing.T) {

	conf, err := config.LoadConfig("../.")
	require.NoError(t, err)

	client := InitRedis(conf)

	err = client.Ping(context.Background()).Err()
	require.NoError(t, err)

	key := "key"
	value := "value"
	err = client.Set(context.Background(), key, value, 0).Err()
	require.NoError(t, err)

	res, err := client.Get(context.Background(), key).Result()
	require.NoError(t, err)
	require.Equal(t, value, res)
}
