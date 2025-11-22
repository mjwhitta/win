//go:build windows

package main

import (
	"fmt"
	"os"

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
		fmt.Println(id.Whoami())
		os.Exit(Good)
	}

	if flags.user {
		fmt.Println(id.WhoamiUser())
		fmt.Println()
	}

	if flags.groups {
		fmt.Println(id.WhoamiGroups())
		fmt.Println()
	}

	if flags.privs {
		fmt.Println(id.WhoamiPriv())
		fmt.Println()
	}
}
