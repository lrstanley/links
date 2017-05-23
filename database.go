package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/boltdb/bolt"
	"github.com/timshannon/bolthold"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func uuid(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Link represents a url that we are shortening.
type Link struct {
	UID            string
	URL            string    // The URL we're expanding.
	Created        time.Time // When the link was submitted.
	Hits           int       // How many times we've expanded for users.
	Author         string    // IP address of request. May be blank.
	EncryptionHash string    // Used to password protect (sha256).
}

func (l *Link) AddHit() {
	l.Hits++

	db := newDB(false)
	if err := db.Upsert(l.UID, l); err != nil {
		debug.Printf("unable to increment hits on %s: %s", l.UID, err)
	}
	db.Close()
}

func (l *Link) Short() string {
	return conf.Site + "/" + l.UID
}

func (l *Link) CheckHash(input string) bool {
	return hash(input) == l.EncryptionHash
}

func hash(input string) string {
	return fmt.Sprintf("%s", sha256.Sum256([]byte(input)))
}

func newDB(readOnly bool) *bolthold.Store {
	store, err := bolthold.Open(conf.DBPath, 0660, &bolthold.Options{Options: &bolt.Options{
		ReadOnly: readOnly,
	}})
	if err != nil {
		panic(fmt.Sprintf("unable to open db: %s", err))
	}

	return store
}

func verifyDB() {
	debug.Printf("verifying access to db: %s", conf.DBPath)
	db := newDB(false)
	db.Close()
	debug.Print("successfully verified access to db")
}
