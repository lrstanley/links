// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
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
}

// Link represents a url that we are shortening.
type Link struct {
	UID            string    `boltholdKey`
	URL            string    // The URL we're expanding.
	Created        time.Time // When the link was submitted.
	Hits           int       // How many times we've expanded for users.
	Author         string    // IP address of request. May be blank.
	EncryptionHash string    // Used to password protect (sha256).
}

func (l *Link) Create() error {
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

	db := newDB(false)
	defer db.Close()

	// Check for dups.
	var result []Link
	err = db.Find(&result, bolthold.Where("URL").Eq(l.URL).And("EncryptionHash").Eq(l.EncryptionHash).Limit(1))
	if err != nil {
		panic(err)
	}

	// Assume there is a dup, just return it to the user.
	if len(result) > 0 {
		l.UID = result[0].UID
		return nil
	}

	// Store it.
	for {
		l.UID = uuid(4)
		err = db.Insert(l.UID, l)
		if err != nil {
			if err == bolthold.ErrKeyExists {
				// Keep looping through until we're able to store one which
				// doesn't collide with a pre-existing key.
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
	store, err := bolthold.Open(conf.DBPath, 0660, &bolthold.Options{Options: &bolt.Options{
		ReadOnly: readOnly,
		Timeout:  10 * time.Second,
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

func migrateDB(mysqlAuth string) {
	sdb, err := sql.Open("mysql", mysqlAuth)
	if err != nil {
		debug.Fatalf("unable to open sql connection: %s", err)
	}
	defer sdb.Close()

	if err = sdb.Ping(); err != nil {
		debug.Fatalf("unable to ping mysql server: %s", err)
	}

	var rows *sql.Rows
	rows, err = sdb.Query("SELECT id, url, hits, COALESCE(password, '') as password, time FROM urls")
	if err != nil {
		debug.Fatal(err)
	}
	defer rows.Close()

	db := newDB(false)
	defer db.Close()

	var count, hits, errors int

	for rows.Next() {
		link := Link{}
		var ts []uint8
		if err = rows.Scan(&link.UID, &link.URL, &link.Hits, &link.EncryptionHash, &ts); err != nil {
			debug.Printf("error scanning row: %s", err)
			errors++
			continue
		}

		tsNano, _ := strconv.ParseInt(string(ts), 10, 64)
		link.Created = time.Unix(tsNano, 0)

		hits += link.Hits

		err = db.Insert(link.UID, &link)
		if err != nil {
			debug.Printf("error inserting exported link into db: %s", err)
			errors++
			continue
		}

		count++
		debug.Printf("count[%4d]: importing %s -> %s", count, link.UID, link.URL)
	}
	if err = rows.Err(); err != nil {
		debug.Fatal(err)
	}

	incGlobalStats(db, hits, count)

	debug.Print("complete")
	debug.Printf("%d imported, %d errors", count, errors)
}
