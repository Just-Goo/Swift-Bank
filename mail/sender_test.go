package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zde37/Swift_Bank/config"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := config.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSender, config.EmailAddress, config.EmailPassword)

	subject := "A test email"
	content := `
	<h1>Hello World</h1>
	<p>This is a test message from <a href="https://github.com/zde37">ZDE</a></p>
	`
	to := []string{""} // provide an email address to send to
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
