// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/timshannon/bolthold"
	bolt "go.etcd.io/bbolt"
)

const collisionMax = 20

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func uuid(n int64) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type GlobalStats struct {
	mu        sync.RWMutex
	Shortened int
	Redirects int
}

func incGlobalStats(db *bolthold.Store, hits, created int) {
	tmp := GlobalStats{}
	err := db.Get("global_stats", &tmp)

	if err != nil && err != bolthold.ErrNotFound {
		debug.Printf("unable to read global stats: %s", err)
		return
	}

	tmp.Redirects += hits
	tmp.Shortened += created

	err = db.Upsert("global_stats", &tmp)
	if err != nil {
		debug.Printf("unable to save global stats: %s", err)
	}

	updateGlobalStats(db)
}

var cachedGlobalStats = GlobalStats{}

func updateGlobalStats(db *bolthold.Store) {
	if db == nil {
		db = newDB(true)
		defer db.Close()
	}

	tmp := GlobalStats{}

	err := db.Get("global_stats", &tmp)

	if err != nil && err != bolthold.ErrNotFound {
		debug.Printf("unable to read global stats: %s", err)
		return
	}

	cachedGlobalStats.mu.Lock()
	cachedGlobalStats.Redirects = tmp.Redirects
	cachedGlobalStats.Shortened = tmp.Shortened
	cachedGlobalStats.mu.Unlock()

	metricRedirected.Set(float64(tmp.Redirects))
	metricShortened.Set(float64(tmp.Shortened))
}

// Link represents a url that we are shortening.
type Link struct {
	UID            string    `boltholdKey:"UID"`
	URL            string    // The URL we're expanding.
	Created        time.Time // When the link was submitted.
	Hits           int       // How many times we've expanded for users.
	Author         string    // IP address of request. May be blank.
	EncryptionHash string    // Used to password protect (sha256).
}

func (l *Link) Create(db *bolthold.Store) error {
	if len(l.URL) < 1 {
		return errors.New("please supply a url to shorten")
	}

	uri, err := url.Parse(l.URL)
	if err != nil || uri.Hostname() == "" {
		return errors.New("unable to parse url: " + l.URL)
	}

	if !isValidScheme(uri.Scheme) {
		return errors.New("invalid url scheme. allowed schemes: " + strings.Join(validSchemes, ", "))
	}

	if strings.Contains(strings.ToLower(conf.Site), strings.ToLower(uri.Hostname())) {
		return errors.New("can't shorten a link for " + conf.Site)
	}

	l.URL = uri.String()
	l.Created = time.Now()

	if db == nil {
		db = newDB(false)
		defer db.Close()
	}

	// Store it. If the UID was pre-defined, don't generate one.
	if l.UID != "" {
		return db.Insert(l.UID, l)
	}

	var collisionCount int
	for {
		l.UID = uuid(atomic.LoadInt64(&conf.KeyLength))
		err = db.Insert(l.UID, l)
		if err != nil {
			if err == bolthold.ErrKeyExists {
				// Keep looping through until we're able to store one which
				// doesn't collide with a pre-existing key.
				collisionCount++
				if collisionCount >= collisionMax {
					// If we continue to get collisions in a row, bump the global
					// key value to one higher. Worst case if the collision max is
					// hit, it will only be computationally difficult on the first
					// shortened link upon restart.
					collisionCount = 0
					newKeyLength := atomic.AddInt64(&conf.KeyLength, 1)
					debug.Printf("collision maximum hit (count: %d); increasing global key length from %d to %d", collisionMax, newKeyLength-1, newKeyLength)
				}
				continue
			}

			panic(err)
		}

		break
	}

	incGlobalStats(db, 0, 1)

	return nil
}

func (l *Link) AddHit() {
	l.Hits++

	db := newDB(false)
	if err := db.Update(l.UID, l); err != nil {
		debug.Printf("unable to increment hits on %s: %s", l.UID, err)
	}
	incGlobalStats(db, 1, 0)
	db.Close()
}

func (l *Link) Short() string {
	return conf.Site + "/" + l.UID
}

func (l *Link) CheckHash(input string) bool {
	return hash(input) == l.EncryptionHash
}

func hash(input string) string {
	if input == "" {
		return ""
	}

	out := sha256.Sum256([]byte(input))

	return hex.EncodeToString(out[:])
}

func newDB(readOnly bool) *bolthold.Store {
	store, err := bolthold.Open(conf.DBPath, 0o660, &bolthold.Options{Options: &bolt.Options{
		FreelistType: bolt.FreelistMapType,
		ReadOnly:     readOnly,
		Timeout:      25 * time.Second,
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

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func dbExportJSON(path string) {
	f, err := os.Create(path)
	if err != nil {
		debug.Fatalf("error occurred while trying to open %q: %s", path, err)
	}
	defer f.Close()

	gob.Register(Link{})

	db := newDB(true)
	defer db.Close()

	var n int

	db.Bolt().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(getType(Link{})))

		b.ForEach(func(k, v []byte) error {
			l := Link{}

			dec := gob.NewDecoder(bytes.NewBuffer(v))
			err = dec.Decode(&l)
			if err != nil {
				debug.Printf("failure: decode %s: %s", k, err)
				return err
			}

			out, _ := json.Marshal(&l)
			out = append(out, 0x0a)
			n, err = f.Write(out)
			if err != nil {
				debug.Printf("failure: unable to write %s to %s: %s", k, path, err)
				return err
			}

			debug.Printf("success: exported %s (%d bytes)", l.Short(), n)
			return nil
		})
		return nil
	})
}
