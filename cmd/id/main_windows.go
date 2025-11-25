//go:build windows

package main

import (
	"fmt"

	"github.com/mjwhitta/log"
	"github.com/mjwhitta/win/user"
)

var id *user.ID

func init() {
	var e error

	if id, e = user.Identity(); e != nil {
		panic(e)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	validate()

	if !flags.groups && !flags.privs && !flags.user {
		fmt.Println(id.Whoami())
		return
	}

	if flags.user {
		fmt.Println(id.WhoamiUser())
	}

	if flags.groups {
		if flags.user {
			fmt.Println()
		}

		fmt.Println(id.WhoamiGroups())
	}

	if flags.privs {
		if flags.groups || flags.user {
			fmt.Println()
		}

		fmt.Println(id.WhoamiPriv())
	}

	if flags.all {
		fmt.Println()
	}
}
