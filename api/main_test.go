package api

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	u "backend-masterclass/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store sqlc.Store) (server *Server) {
	// Setting up a fake config.
	cfg := u.Config{
		TokenKey:            u.RandomStr(32),
		AccessTokenDuration: time.Minute,
	}

	pasetoTokenMaker, err := token.NewPasetoMaker(cfg.TokenKey)
	require.NoError(t, err)

	server = NewServer(store, cfg, pasetoTokenMaker)

	return
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())

}
