package crypt

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/billybobjoeaglt/chatlab/config"

	"golang.org/x/crypto/openpgp"
)

var privateKeyEntityList openpgp.EntityList
var passphrase string
var keyMap = make(map[string]*openpgp.Entity)

func getKeyByKeybaseUsername(username string) (openpgp.EntityList, error) {
	// Gets public key of recipient
	resp, err := http.Get("https://keybase.io/" + username + "/key.asc")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Reads key into list of keys
	entityList, err := openpgp.ReadArmoredKeyRing(resp.Body)
	if err != nil {
		return nil, err
	}
	return entityList, nil
}

func Encrypt(message string, users []string) (string, error) {

	var entityList openpgp.EntityList

	for _, username := range users {
		val, ok := keyMap[username]
		if !ok {
			eL, err := getKeyByKeybaseUsername(username)
			if err != nil {
				return "", err
			}
			keyMap[username] = eL[0]
			val, _ = keyMap[username]
		}
		entityList = append(entityList, val)
	}

	// New buffer where the result of the encripted msg will be
	buf := new(bytes.Buffer)

	if privateKeyEntityList == nil {
		createPrivKey()
	}

	// Create an encryption stream
	plaintext, err := openpgp.Encrypt(buf, entityList, privateKeyEntityList[0], nil, nil)
	if err != nil {
		return "", err
	}

	// Write a byte array saying the message to encryption stream

	if _, err := plaintext.Write([]byte(message)); err != nil {
		return "", err
	}

	// Close streams, finishing encryption and armor texts
	plaintext.Close()

	base64Enc := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Enc, nil
}

func createPrivKey() error {
	var err error

	if passphrase == "" {
		var pass []byte
		pass, err = ioutil.ReadFile(config.GetConfig().Passphrase)
		if err != nil {
			return err
		}
		passphrase = strings.TrimSpace(string(pass))
	}
	var keyringFileBuffer *os.File
	keyringFileBuffer, err = os.Open(config.GetConfig().PrivateKey)
	if err != nil {
		return err
	}
	defer keyringFileBuffer.Close()

	privateKeyEntityList, err = openpgp.ReadArmoredKeyRing(keyringFileBuffer)

	entity := privateKeyEntityList[0]
	passphraseByte := []byte(passphrase)
	if entity.PrivateKey.Encrypted {
		if err = entity.PrivateKey.Decrypt(passphraseByte); err != nil {
			return err
		}
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey.Encrypted {
			if err = subkey.PrivateKey.Decrypt(passphraseByte); err != nil {
				return err
			}
		}
	}

	return err
}

func Decrypt(base64msg string) (*openpgp.MessageDetails, error) {

	var err error

	messageByteArr, err := base64.StdEncoding.DecodeString(base64msg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(messageByteArr)

	if privateKeyEntityList == nil {
		createPrivKey()
	}

	md, err := openpgp.ReadMessage(buf, privateKeyEntityList, nil, nil)
	if err != nil {
		return nil, err
	}
	return md, nil
}
