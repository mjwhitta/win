//go:build windows

package main

import (
	"fmt"

	"github.com/mjwhitta/log"
	w32 "github.com/mjwhitta/win/api"
	"github.com/mjwhitta/win/proc"
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
	var procs []*w32.ProcessEntry32

	validate()

	if procs, e = proc.Get(); e != nil {
		panic(e)
	}

	fmt.Printf("PPID\tPID\tTHREADS\tNAME\n")
	for _, pe := range procs {
		fmt.Printf(
			"%d\t%d\t%d\t%s\n",
			pe.ParentPID,
			pe.PID,
			pe.ThreadCount,
			pe.ExeFile(),
		)
	}
}
