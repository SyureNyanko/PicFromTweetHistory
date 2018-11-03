package main

import (
"testing"
)

var (
	testjspath = "./testcontents/2018_10.js"
	downloadpath = "./testcontents"
)


func TestPicDownloader(t *testing.T) {
	PicDownloader(testjspath, downloadpath)
}