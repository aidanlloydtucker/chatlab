package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
)

// XXX: For actual encryption in app do not use armor for decoding

var privateKeyEntityList openpgp.EntityList

func decrypt(message string, passphrase string) (string, error) {

	var err error

	buf := bytes.NewBuffer([]byte(message))

	/*result, err := armor.Decode(buf)
	if err != nil {
		return "", err
	}*/
	if privateKeyEntityList == nil {
		var keyringFileBuffer *os.File
		keyringFileBuffer, err = os.Open(config.PrivateKey)
		if err != nil {
			return "", err
		}
		defer keyringFileBuffer.Close()

		privateKeyEntityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
		if err != nil {
			return "", err
		}
	}

	entity := privateKeyEntityList[0]
	passphraseByte := []byte(passphrase)
	if entity.PrivateKey.Encrypted {
		if err = entity.PrivateKey.Decrypt(passphraseByte); err != nil {
			return "", err
		}
	}
	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey.Encrypted {
			if err = subkey.PrivateKey.Decrypt(passphraseByte); err != nil {
				return "", err
			}
		}
	}

	md, err := openpgp.ReadMessage(buf, privateKeyEntityList, nil, nil)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

/*func main() {
	msg, err := ioutil.ReadFile("./msg.txt")
	if err != nil {
		panic(err)
	}

	pass, err := ioutil.ReadFile("./pass.key")
	if err != nil {
		panic(err)
	}

	str, err := decrypt(string(msg), strings.TrimSpace(string(pass)))
	if err != nil {
		panic(err)
	}

	fmt.Println(str)
}*/
