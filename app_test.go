package main

import "testing"

func Test(t *testing.T) {
	erg := ReturnDouble(5)

	if erg != 10 {
		t.Error("test failed")
	}
}
