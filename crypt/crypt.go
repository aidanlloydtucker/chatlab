package crypt

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"net/http"
	"os"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"

	"golang.org/x/crypto/openpgp"
)

var privateKeyEntityList openpgp.EntityList
var keyMap = make(map[string]*openpgp.Entity)

func GetKeyMap() map[string]*openpgp.Entity {
	return keyMap
}

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

func AddPublicKeyToMap(username string, key string) {
	eL, err := openpgp.ReadKeyRing(bytes.NewReader([]byte(key)))
	if err != nil {
		return
	}
	keyMap[username] = eL[0]
}

// Encrypts a message to the users
func Encrypt(message string, users []string) (string, error) {
	var entityList openpgp.EntityList

	// Loop through and get the key for each of the users
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

// Encrypts a message struct using the gob protocol
func EncryptMessage(msg common.Message) (string, error) {
	var gobBuf bytes.Buffer
	enc := gob.NewEncoder(&gobBuf) // Will write to buf.
	err := enc.Encode(msg)
	if err != nil {
		return "", err
	}

	var entityList openpgp.EntityList

	// Loop through and get the key for each of the users
	for _, username := range msg.ToUsers {
		val, ok := keyMap[username]
		if !ok {
			eL, err2 := getKeyByKeybaseUsername(username)
			if err2 != nil {
				return "", err2
			}
			keyMap[username] = eL[0]
			val, _ = keyMap[username]
		}
		entityList = append(entityList, val)
	}

	// New buffer where the result of the encripted msg will be

	if privateKeyEntityList == nil {
		createPrivKey()
	}

	encBuf := new(bytes.Buffer)

	// Create an encryption stream
	plaintext, err := openpgp.Encrypt(encBuf, entityList, privateKeyEntityList[0], nil, nil)
	if err != nil {
		return "", err
	}

	// Write a byte array saying the message to encryption stream

	if _, err := plaintext.Write(gobBuf.Bytes()); err != nil {
		return "", err
	}

	// Close streams, finishing encryption and armor texts
	plaintext.Close()

	base64Enc := base64.StdEncoding.EncodeToString(encBuf.Bytes())

	return base64Enc, nil
}

func createPrivKey() error {
	var err error

	// Read private key from disk
	var keyringFileBuffer *os.File
	keyringFileBuffer, err = os.Open(config.GetConfig().PrivateKey)
	if err != nil {
		return err
	}
	defer keyringFileBuffer.Close()

	// Read key int var
	privateKeyEntityList, err = openpgp.ReadArmoredKeyRing(keyringFileBuffer)

	entity := privateKeyEntityList[0]
	passphraseByte := []byte(config.Password)

	// Decrypt private key with password
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

// Decrypts an unarmored base64 message
func Decrypt(base64msg string) (*openpgp.MessageDetails, error) {

	var err error

	// Turn message into []byte
	messageByteArr, err := base64.StdEncoding.DecodeString(base64msg)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(messageByteArr)

	if privateKeyEntityList == nil {
		createPrivKey()
	}

	// Decrypts message with private key and buffer of []byte of message
	md, err := openpgp.ReadMessage(buf, privateKeyEntityList, nil, nil)
	if err != nil {
		return nil, err
	}
	return md, nil
}

// Decrypts an EncryptedMessage using the gob protocol
func DecryptMessage(base64msg string) (*openpgp.MessageDetails, *common.Message, error) {
	var err error

	// Turn message into []byte
	messageByteArr, err := base64.StdEncoding.DecodeString(base64msg)
	if err != nil {
		return nil, nil, err
	}

	decBuf := bytes.NewBuffer(messageByteArr)

	if privateKeyEntityList == nil {
		createPrivKey()
	}

	// Decrypts message with private key and buffer of []byte of message
	md, err := openpgp.ReadMessage(decBuf, privateKeyEntityList, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	dec := gob.NewDecoder(md.UnverifiedBody)

	msg := common.Message{}
	err = dec.Decode(&msg)
	if err != nil {
		return nil, nil, err
	}

	return md, &msg, nil

}
