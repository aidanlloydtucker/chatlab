package main

func main() {
	// Passphrase for private key
	/*var passprase, err = ioutil.ReadFile("./pass.key")
	if err != nil {
		panic(err)
	}*/

	go printAll(outputChannel)
	listen()
}
