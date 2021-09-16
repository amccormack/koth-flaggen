package main

import (
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"sync"
	"net/http"
	"encoding/json"
)

var seed string = "defaultsecret"
var checkPath string = "/secret.txt"
var fc *FlagCheck

type FlagCheck struct {
	mu sync.Mutex
	flags map[string]bool
	points map[string]map[string]bool
}
func NewFlagCheck() *FlagCheck {
	fc := &FlagCheck{
		flags: make(map[string]bool),
		points: make(map[string]map[string]bool),
	}
	return fc
}

func (fc *FlagCheck) SetFlags(t time.Time) {
	f1, f2, f3 := GenFlags(t)
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.flags = make(map[string]bool)
	fc.flags[f1] = true
	fc.flags[f2] = true
	fc.flags[f3] = true

	fmt.Printf("SF %v: %v %v %v\n", t.Format("03:04:05"), f1, f2, f3)
}

func (fc *FlagCheck) CheckFlag(guess string) bool {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	_, exists := fc.flags[guess]
	return exists
}

func (fc *FlagCheck) RecordFlag(guess string, user string) bool {
	if !fc.CheckFlag(guess){
		return false
	}
	_, exists := fc.points[user]
	if !exists {
		userMap := make(map[string]bool)
		fc.points[user] = userMap
	}
	fc.points[user][guess] = true
	return true
}

func (fc *FlagCheck) GetPoints() []UserPoint {
	points := make([]UserPoint, len(fc.points))
	index :=0
	for user, flags := range fc.points {
		points[index].Username = user
		points[index].Points = len(flags)
		index++
	}
	return points
}

type UserPoint struct {
	Username string
	Points int
}

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

func GenFlags(current time.Time) (string, string, string) {
	prev := current.Add(-1 * time.Minute)
	next := current.Add(time.Minute)

	ps := prev.Format("2006-01-02 15:04")
	cs := current.Format("2006-01-02 15:04")
	ns := next.Format("2006-01-02 15:04")

	return ComputeHmac256(ps, seed), ComputeHmac256(cs, seed), ComputeHmac256(ns, seed)
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/score" {
		scores := fc.GetPoints()
		msg, err := json.Marshal(scores)
		if err == nil {
			fmt.Fprintf(w, "%s\n", msg)	
		} else {
			http.Error(w, "Error encoding...", http.StatusInternalServerError)
		}
		return
	}

	q := r.URL.Query()
	flag := q.Get("flag")
	user := q.Get("user")
	if user == "" || len(flag) != 64 {
		message := fmt.Sprintf("Invalid user and flag\nuser: '%s' flag: '%s'\n", user, flag)
		message = message + fmt.Sprintf("Submit flag and user like /?flag=...&user=...")
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "user: '%s' flag: '%s'\n", user, flag)
	if fc.RecordFlag(flag, user){
		fmt.Fprintf(w, "yay\n")
	} else {
		fmt.Fprintf(w, "boo\n")
	}
}

func main() {
	fc = NewFlagCheck()
	loadSeed()

	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)
	fc.SetFlags(time.Now())

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fc.SetFlags(t)
			}
		}
	}()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)

	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel.
	time.Sleep(60 * time.Second)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")


}
