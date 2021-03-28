package mode_s

import (
	"testing"
)

func TestIsNoop(t *testing.T) {
	frames := []Frame{
		{raw: ""},
		{raw: "0"},
		{raw: "@00000"},
		{raw: "*0"},
		{raw: "*0000"},
	}
	for _, f := range frames {
		if !f.isNoOp() {
			t.Errorf("Failed to detect NoOp frame: %s", f.raw)
		}
	}
}

func TestIsNotNoop(t *testing.T) {
	frames := []Frame{
		{raw: "10"},
		{raw: "123"},
		{raw: "@123;"},
		{raw: "*3"},
		{raw: "*023"},
		{raw: "*00001"},
	}
	for _, f := range frames {
		f.full = "*" + f.raw + ";"
		if f.isNoOp() {
			t.Errorf("Failed detect non NoOp frame as NoOp: %s", f.raw)
		}
	}
}

func TestFrame_isNoOp(t *testing.T) {
	type fields struct {
		raw string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "noop", fields: fields{raw: ""}, want: true},
		{name: "noop", fields: fields{raw: "0"}, want: true},
		{name: "noop", fields: fields{raw: "00"}, want: true},
		{name: "noop", fields: fields{raw: "000"}, want: true},
		{name: "noop", fields: fields{raw: "0000"}, want: true},
		{name: "noop", fields: fields{raw: "00000"}, want: true},
		{name: "noop", fields: fields{raw: "000000"}, want: true},
		{name: "noop", fields: fields{raw: "0000000"}, want: true},
		{name: "noop", fields: fields{raw: "00000000"}, want: true},
		{name: "noop", fields: fields{raw: "000000000"}, want: true},
		{name: "noop", fields: fields{raw: "0000000000"}, want: true},
		{name: "noop", fields: fields{raw: "00000000000"}, want: true},
		{name: "noop", fields: fields{raw: "000000000000"}, want: true},
		{name: "noop", fields: fields{raw: "0000000000000"}, want: true},
		{name: "noop", fields: fields{raw: "00000000000000"}, want: true},
		{name: "noop", fields: fields{raw: "000000000000000"}, want: true},
		{name: "noop", fields: fields{raw: "0000000000000000"}, want: true},
		{name: "noop", fields: fields{raw: "00000000000000000"}, want: false},
		{name: "bad", fields: fields{raw: "1"}, want: false},
		{name: "bad", fields: fields{raw: "12"}, want: false},
		{name: "bad", fields: fields{raw: "123"}, want: false},
		{name: "bad", fields: fields{raw: "1234"}, want: false},
		{name: "bad", fields: fields{raw: "12345"}, want: false},
		{name: "bad", fields: fields{raw: "123456"}, want: false},
		{name: "bad", fields: fields{raw: "1234567"}, want: false},
		{name: "bad", fields: fields{raw: "12345678"}, want: false},
		{name: "bad", fields: fields{raw: "123456789"}, want: false},
		{name: "bad", fields: fields{raw: "1234567890"}, want: false},
		{name: "bad", fields: fields{raw: "12345678901"}, want: false},
		{name: "bad", fields: fields{raw: "123456789012"}, want: false},
		{name: "bad", fields: fields{raw: "1234567890123"}, want: false},
		{name: "bad", fields: fields{raw: "12345678901234"}, want: false},
		{name: "bad", fields: fields{raw: "123456789012345"}, want: false},
		{name: "bad", fields: fields{raw: "1234567890123456"}, want: false},
		{name: "bad", fields: fields{raw: "12345678901234567"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Frame{
				raw:  tt.fields.raw,
				full: "*" + tt.fields.raw + ";",
			}
			if got := f.isNoOp(); got != tt.want {
				t.Errorf("for %s isNoOp() = %v, want %v", tt.fields.raw, got, tt.want)
			}
		})
	}
}
