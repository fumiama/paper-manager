package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/fumiama/paper-manager/backend/api"
	"github.com/fumiama/paper-manager/backend/file"
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
		logrus.Errorln("[net.Listen]\t", err)
		os.Exit(line())
	}

	http.HandleFunc("/api/", api.Handler)
	http.HandleFunc("/file/", file.Handler)

	logrus.Infoln("[http.Serve]\t start at", l.Addr())
	logrus.Errorln("[http.Serve]\t", http.Serve(l, nil))
	os.Exit(line())
}
