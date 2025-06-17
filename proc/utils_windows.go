//go:build windows

package proc

import w32 "github.com/mjwhitta/win/api"

func psLessFunc(procs []*w32.ProcessEntry32) func(i int, j int) bool {
	return func(i int, j int) bool {
		if procs[i].ParentPID == procs[j].ParentPID {
			return procs[i].PID < procs[j].PID
		}

		return procs[i].ParentPID < procs[j].ParentPID
	}
}
