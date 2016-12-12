package yolp

import "testing"

func TestError(t *testing.T) {
	Test{"Failed", Error{Message: "Failed"}.Error()}.Compare(t)
}
