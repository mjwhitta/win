package api

// WIN32_LEAN_AND_MEAN...
// and accctrl + stdlib + tlhelp32 + winapifamily + winhttp + wininet
// + winnt + winuser
//go:generate go run ../tools/defines.go api accctrl.h cderr.h dde.h ddeml.h dlgs.h lzexpand.h mmsystem.h nb30.h rpc.h shellapi.h stdlib.h tlhelp32.h winapifamily.h wincrypt.h winefs.h winhttp.h wininet.h winnt.h winperf.h winscard.h winsock.h winuser.h
