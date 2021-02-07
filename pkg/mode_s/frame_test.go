package mode_s

import "testing"

func TestIsNoop(t *testing.T) {
	frames := []Frame{
		{
			raw: "",
		},
		{
			raw: "0",
		},
		{
			raw: "@00000",
		},
		{
			raw: "*0",
		},
		{
			raw: "*0000",
		},
	}
	for _, f := range frames {
		if !f.isNoOp() {
			t.Errorf("Failed to detect NoOp frame: %s", f.raw)
		}
	}
}

func TestIsNotNoop(t *testing.T) {
	frames := []Frame{
		{
			raw: "10",
		},
		{
			raw: "123",
		},
		{
			raw: "@123;",
		},
		{
			raw: "*3",
		},
		{
			raw: "*023",
		},
		{
			raw: "*00001",
		},
	}
	for _, f := range frames {
		if f.isNoOp() {
			t.Errorf("Failed detect non NoOp frame as NoOp: %s", f.raw)
		}
	}
}
