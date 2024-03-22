//go:build windows

package main

import (
	"os"

	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/win/user"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r.(error).Error())
			}
			log.ErrX(Exception, r.(error).Error())
		}
	}()

	var e error
	var id *user.ID

	validate()

	if id, e = user.Identity(); e != nil {
		panic(e)
	}

	if !(flags.groups || flags.privs || flags.user) {
		hl.Println(id.Whoami())
		os.Exit(Good)
	}

	if flags.user {
		hl.Println(id.WhoamiUser())
		hl.Println()
	}

	if flags.groups {
		hl.Println(id.WhoamiGroups())
		hl.Println()
	}

	if flags.privs {
		hl.Println(id.WhoamiPriv())
		hl.Println()
	}
}
