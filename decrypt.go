package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
    "io/ioutil"
    "os"
    "strings"
)

// XXX: For actual encryption in app do not use armor for decoding

func decrypt(message string, passphrase string) (string, error) {

	buf := bytes.NewBuffer([]byte(message))

	result, err := armor.Decode(buf)
	if err != nil {
		return "", err
	}

    keyringFileBuffer, err := os.Open("./key.key")
    if err != nil {
		return "", err
	}
    defer keyringFileBuffer.Close()

    entityList, err := openpgp.ReadArmoredKeyRing(keyringFileBuffer)
    if err != nil {
        return "", err
    }

    entity := entityList[0]
    passphraseByte := []byte(passphrase)
    if entity.PrivateKey.Encrypted {
        err := entity.PrivateKey.Decrypt(passphraseByte)
        if err != nil {
            return "", err
        }
    }
    for _, subkey := range entity.Subkeys {
        if subkey.PrivateKey.Encrypted {
            err := subkey.PrivateKey.Decrypt(passphraseByte)
            if err != nil {
                return "", err
            }
        }
    }

	md, err := openpgp.ReadMessage(result.Body, entityList, nil, nil)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
    if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func main() {
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
}
