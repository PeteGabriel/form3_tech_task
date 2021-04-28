package form3_task

import "testing"

func TestForm3(t *testing.T) {
	if !Tech3() {
		t.Fail()
	}
}