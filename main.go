package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/fumiama/paper-manager/backend/api"
)

func line() int {
	_, _, fileLine, ok := runtime.Caller(1)
	if ok {
		return fileLine
	}
	return -1
}

func main() {
	addr := flag.String("l", "[::]:3000", "listen addr")
	flag.Parse()
	l, err := net.Listen("tcp", *addr)
	if err != nil {
		logrus.Errorln("[net.Listen]", err)
		os.Exit(line())
	}

	http.HandleFunc("/api/", api.Handler)
	http.HandleFunc("/file/", api.FileHandler)
	http.HandleFunc("/upload", api.UploadHandler)

	logrus.Infoln("[http.Serve] start at", l.Addr())
	logrus.Errorln("[http.Serve]", http.Serve(l, nil))
	os.Exit(line())
}
