package bot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wneessen/go-wftk/crypto/random"
	"os"
	"os/user"
)

func (b *Bot) GetEncryptionKey() error {
	l := log.WithFields(log.Fields{
		"action": "bot.GetEncryptionKey",
	})
	userObj, err := user.Current()
	if err != nil {
		return err
	}
	userHome := userObj.HomeDir

	encKeyFile := fmt.Sprintf("%v/.sotbot_enc_key", userHome)
	encKeyFileExists := true
	l.Debugf("Checking if %v exists...", encKeyFile)
	_, err = os.Stat(encKeyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		encKeyFileExists = false
	}

	if !encKeyFileExists {
		l.Debugf("Key file does not exist. Creating it...")
		encKey, err := random.GenerateRandomString(32, false, false)
		if err != nil {
			return err
		}
		if err := os.WriteFile(encKeyFile, []byte(encKey), 0600); err != nil {
			return err
		}
		b.Config.Set("enc_key", encKey)
		return nil
	}

	l.Debugf("Key file already exists. Reading it...")
	encKeyBytes, err := os.ReadFile(encKeyFile)
	if err != nil {
		return err
	}
	b.Config.Set("enc_key", string(encKeyBytes))
	return nil
}
