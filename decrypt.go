package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
)

// XXX: For actual encryption in app do not use armor for decoding

var privateKeyEntityList openpgp.EntityList
var passphrase string

func createPrivKey() error {
	var err error

	if passphrase == "" {
		var pass []byte
		pass, err = ioutil.ReadFile(config.Passphrase)
		if err != nil {
			return err
		}
		passphrase = strings.TrimSpace(string(pass))
	}
	var keyringFileBuffer *os.File
	keyringFileBuffer, err = os.Open(config.PrivateKey)
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

func decrypt(base64msg string) (string, error) {

	var err error

	messageByteArr, err := base64.StdEncoding.DecodeString(base64msg)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(messageByteArr)

	/*result, err := armor.Decode(buf)
	if err != nil {
		return "", err
	}*/
	if privateKeyEntityList == nil {
		createPrivKey()
	}

	md, err := openpgp.ReadMessage(buf, privateKeyEntityList, nil, nil)
	if err != nil {
		return "", err
	}
	/*for k := range md.SignedBy.Entity.Identities {
		fmt.Println("Name: " + md.SignedBy.Entity.Identities[k].UserId.Name)
		fmt.Println("Email: " + md.SignedBy.Entity.Identities[k].UserId.Email)
		fmt.Println("Comment: " + md.SignedBy.Entity.Identities[k].UserId.Comment)
		fmt.Println("Creation Time: " + md.SignedBy.Entity.Identities[k].SelfSignature.CreationTime.Format(time.UnixDate) + "\n")
	}*/

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
