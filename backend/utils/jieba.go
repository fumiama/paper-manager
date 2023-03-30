package utils

import (
	"archive/zip"
	"bytes"
	_ "embed"

	"github.com/fumiama/jieba"
)

//go:embed dict.zip
var dictzip []byte

// Segmenter jieba 分词器
var Segmenter = func() *jieba.Segmenter {
	r, err := zip.NewReader(bytes.NewReader(dictzip), int64(len(dictzip)))
	if err != nil {
		panic(err)
	}
	f, err := r.Open("dict.txt")
	if err != nil {
		panic(err)
	}
	seg, err := jieba.LoadDictionary(f)
	if err != nil {
		panic(err)
	}
	return seg
}()
