package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"

	"github.com/fumiama/paper-manager/backend"
	"github.com/fumiama/paper-manager/frontend"
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

	http.HandleFunc("/api/", backend.APIHandler)
	http.HandleFunc("/file/", backend.FileHandler)
	http.HandleFunc("/paper/", backend.PaperHandler)
	http.HandleFunc("/upload", backend.UploadHandler)
	http.Handle("/", frontend.StaticHandler)

	logrus.Infoln("[http.Serve] start at", l.Addr())
	logrus.Errorln("[http.Serve]", http.Serve(l, nil))
	os.Exit(line())
}
