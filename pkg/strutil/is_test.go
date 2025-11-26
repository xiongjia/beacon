package strutil

import (
	"testing"
)

type MyString string

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"empty string", "", true},
		{"space string", " ", false},
		{"non-empty", "abc", false},
		{"empty custom", MyString(""), true},
		{"space custom", MyString(" "), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.in.(type) {
			case string:
				if got := IsEmpty(v); got != tt.want {
					t.Fatalf("IsEmpty(%q) = %v, want %v", v, got, tt.want)
				}
			case MyString:
				if got := IsEmpty(v); got != tt.want {
					t.Fatalf("IsEmpty(MyString(%q)) = %v, want %v", string(v), got, tt.want)
				}
			default:
				t.Fatalf("unsupported input type %T", tt.in)
			}
		})
	}
}

func TestIsEmptyPtr(t *testing.T) {
	s1 := ""
	s2 := "abc"
	m1 := MyString("")
	m2 := MyString(" ")

	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"nil *string", (*string)(nil), true},
		{"ptr empty string", &s1, true},
		{"ptr non-empty string", &s2, false},
		{"nil *MyString", (*MyString)(nil), true},
		{"ptr empty MyString", &m1, true},
		{"ptr space MyString", &m2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.in.(type) {
			case *string:
				if got := IsEmptyPtr(v); got != tt.want {
					t.Fatalf("IsEmptyPtr(%v) = %v, want %v", v, got, tt.want)
				}
			case *MyString:
				if got := IsEmptyPtr(v); got != tt.want {
					t.Fatalf("IsEmptyPtr(%v) = %v, want %v", v, got, tt.want)
				}
			default:
				t.Fatalf("unsupported input type %T", tt.in)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"empty", "", true},
		{"spaces", "   ", true},
		{"tabs newlines", "\t\n", true},
		{"leading space", " a ", false},
		{"word", "abc", false},
		{"empty custom", MyString(""), true},
		{"spaces custom", MyString("  \t"), true},
		{"non-blank custom", MyString("x"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.in.(type) {
			case string:
				if got := IsBlank(v); got != tt.want {
					t.Fatalf("IsBlank(%q) = %v, want %v", v, got, tt.want)
				}
			case MyString:
				if got := IsBlank(v); got != tt.want {
					t.Fatalf("IsBlank(MyString(%q)) = %v, want %v", string(v), got, tt.want)
				}
			default:
				t.Fatalf("unsupported input type %T", tt.in)
			}
		})
	}
}

func TestIsBlankPtr(t *testing.T) {
	s1 := "   "
	s2 := "a"
	var pnil *string
	m1 := MyString("\t\n")
	m2 := MyString("b")

	tests := []struct {
		name string
		in   any
		want bool
	}{
		{"nil *string", pnil, true},
		{"ptr spaces string", &s1, true},
		{"ptr non-blank string", &s2, false},
		{"nil *MyString", (*MyString)(nil), true},
		{"ptr blank MyString", &m1, true},
		{"ptr non-blank MyString", &m2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.in.(type) {
			case *string:
				if got := IsBlankPtr(v); got != tt.want {
					t.Fatalf("IsBlankPtr(%v) = %v, want %v", v, got, tt.want)
				}
			case *MyString:
				if got := IsBlankPtr(v); got != tt.want {
					t.Fatalf("IsBlankPtr(%v) = %v, want %v", v, got, tt.want)
				}
			default:
				t.Fatalf("unsupported input type %T", tt.in)
			}
		})
	}
}
