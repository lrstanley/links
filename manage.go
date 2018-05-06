// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	sempool "github.com/lrstanley/go-sempool"
	"github.com/timshannon/bolthold"
)

var reCustomID = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,25}$`)

type CommandAdd struct {
	ID string `short:"i" long:"id" description:"custom id to use for shortened link (a-z, A-Z, 0-9, '-' and '_' allowed, 3-25 chars)"`
}

func (*CommandAdd) Usage() string { return "<link> [link]..." }

func (cli *CommandAdd) Execute(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid usage: see 'add --help'")
	}

	if cli.ID != "" {
		if len(args) > 1 {
			return errors.New("invalid usage: can only specify one link to shortened with '--id'")
		}

		if !reCustomID.MatchString(cli.ID) {
			return fmt.Errorf("invalid custom id specified %q: see 'add --help'", cli.ID)
		}
	}

	db := newDB(false)
	defer db.Close()

	pool := sempool.New(3)

	for i := 0; i < len(args); i++ {
		pool.Slot()
		go func(url string) {
			defer pool.Free()
			link := &Link{
				UID:    cli.ID,
				URL:    url,
				Author: "localhost",
			}

			if err := link.Create(db); err != nil {
				fmt.Fprintf(os.Stderr, "error adding %q: %v\n", url, err)
				if len(args) == 1 {
					os.Exit(1)
				}
			}

			fmt.Println(link.Short())
		}(args[i])
	}

	pool.Wait()
	return nil
}

type CommandDelete struct{}

func (*CommandDelete) Usage() string { return "<id or link> [ids or links]..." }

func (cli *CommandDelete) Execute(args []string) error {
	if len(args) < 1 {
		return errors.New("invalid usage: see 'delete --help'")
	}

	db := newDB(false)
	defer db.Close()

	pool := sempool.New(3)

	for i := 0; i < len(args); i++ {
		pool.Slot()
		go func(param string) {
			defer pool.Free()

			var result []Link

			err := db.Find(
				&result,
				bolthold.Where("URL").Eq(param).Or(
					bolthold.Where("UID").Eq(param).Or(
						bolthold.Where("Author").Eq(param),
					),
				),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error searching for %q: %v\n", param, err)
				if len(args) == 1 {
					os.Exit(1)
				}
				return
			}

			for j := 0; j < len(result); j++ {
				err = db.Delete(result[j].UID, &Link{})
				if err != nil {
					fmt.Fprintf(os.Stderr, "error searching for %q: %v\n", param, err)
					if len(args) == 1 {
						os.Exit(1)
					}
					continue
				}
			}
		}(args[i])
	}

	pool.Wait()
	return nil
}
