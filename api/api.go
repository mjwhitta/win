package api

// WIN32_LEAN_AND_MEAN + winapifamily + winhttp + wininet + winnt + winuser
//go:generate go run ../tools/defines.go api cderr.h dde.h ddeml.h dlgs.h lzexpand.h mmsystem.h nb30.h rpc.h shellapi.h winapifamily.h wincrypt.h winefs.h winhttp.h wininet.h winnt.h winperf.h winscard.h winsock.h winuser.h
