package postgres

import (
	"context"
	"testing"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/stretchr/testify/assert"
)

func loadBaseConfig() *common_utils.BaseConfig {
	cfg, err := common_utils.LoadBaseConfig("../", "test")
	if err != nil {
		common_utils.PanicIfError(err)
	}
	return cfg
}

func TestConnectionPostgres(t *testing.T) {
	ctx := context.Background()
	config := loadBaseConfig()
	connPool, err := NewPgxConn(config)
	assert.Nil(t, err)
	assert.NotNil(t, connPool)

	t.Run("TestPing Postgres Connectiom", func(t *testing.T) {
		err = connPool.Ping(ctx)
		assert.Nil(t, err)
	})

	t.Run("TestClose Postgres Connectiom", func(t *testing.T) {
		connPool.Close()
		err = connPool.Ping(ctx)
		assert.NotNil(t, err)
	})

}
