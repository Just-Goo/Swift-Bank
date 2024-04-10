package handler

import (
	"os"
	"testing"
	"time"

	"github.com/Just-Goo/Swift_Bank/config"
	"github.com/Just-Goo/Swift_Bank/helpers"
	"github.com/gin-gonic/gin"
)

var testConfig config.Config

func TestMain(m *testing.M) {

	testConfig = config.Config{
		TokenSymmetricKey: helpers.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
