package main

import (
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

var seed string = "defaultsecret"
var checkPath string = "/secret.txt"

func loadSeed() string{
	flagtype := "flag"

	if _, err := os.Stat(checkPath); os.IsNotExist(err) {
		fmt.Println("Secret not loaded, flag file does not exist.")
		flagtype = "testflag"
	} else {
		sbytes, err := os.ReadFile(checkPath)
		if err != nil {
			fmt.Println("Secret not loaded, flagfile could not be read.")
			flagtype = "filereaderrorflag"
		} else {
			seed = string(sbytes)
		}
	}
	return flagtype
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	flagtype := loadSeed()

	timeString := time.Now().Format("2006-01-02 15:04")
	flag := ComputeHmac256(timeString, seed)
	fmt.Printf("It is currently %v. The %v is:%v \n", timeString, flagtype, flag)
}
