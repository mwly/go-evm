package main

import "testing"

func Test(t *testing.T) {
	erg := 5

	if erg != 10 {
		t.Error("test failed")
	}
}
