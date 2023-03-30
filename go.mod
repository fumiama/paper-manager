module github.com/fumiama/paper-manager

go 1.20

require (
	github.com/FloatTech/sqlite v1.6.0
	github.com/FloatTech/ttl v0.0.0-20220715042055-15612be72f5b
	github.com/RomiChan/syncx v0.0.0-20221202055724-5f842c53020e
	github.com/corona10/goimagehash v1.1.0
	github.com/fumiama/go-base16384 v1.6.4
	github.com/fumiama/go-docx v0.0.0-20230330141738-34f53a967c03
	github.com/fumiama/imgsz v0.0.2
	github.com/fumiama/jieba v0.0.0-20221203025406-36c17a10b565
	github.com/sirupsen/logrus v1.9.0
	golang.org/x/image v0.6.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	modernc.org/libc v1.21.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.4.0 // indirect
	modernc.org/sqlite v1.20.0 // indirect
)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
