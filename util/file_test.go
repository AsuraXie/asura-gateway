package util

import "testing"

func TestFile(t *testing.T) {
	ParseConfFile("../conf/gateway.conf")
}