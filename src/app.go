package main // hello

import (
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {

	seed := "defaultsecret"
	checkPath := "/secret.txt"
	test := ""

	if _, err := os.Stat(checkPath); os.IsNotExist(err) {
		fmt.Println("Secret not loaded, flag is invalid.")
		test = "test"
	} else {
		sbytes, err := os.ReadFile(checkPath)
		if err != nil {
			fmt.Println("Secret not loaded, flagfile could not be read.")
			test = "filereaderror"
		} else {
			seed = string(sbytes)
		}
	}



	timeString := time.Now().Format("2006-01-02 15:04")
	flag := ComputeHmac256(timeString, seed)
	fmt.Printf("It is currently %v. The %vflag is:%v \n", timeString, test, flag)
}
