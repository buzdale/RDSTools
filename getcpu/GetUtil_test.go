package getcpu

import "testing"

func TestHello(t *testing.T) {
	got := Utilization("sldb - 01")
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
