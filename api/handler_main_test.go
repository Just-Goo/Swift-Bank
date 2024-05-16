package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zde37/Swift_Bank/config"
	"github.com/zde37/Swift_Bank/helpers"
)

var testConfig config.Config

func TestMain(m *testing.M) {

	testConfig = config.Config{
		TokenSymmetricKey:   helpers.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
