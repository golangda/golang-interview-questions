package word

import "testing"

func TestCount(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"  ", 0},
		{"hello", 1},
		{"go is fun", 3},
	}
	for _, c := range cases {
		if got := Count(c.in); got != c.want {
			t.Fatalf("Count(%q) = %d; want %d", c.in, got, c.want)
		}
	}
}

func TestRepeatCount(t *testing.T) {
	got := RepeatCount("go is fun", 10000)
	if got <= 0 {
		t.Fatal("unexpected result")
	}
}

func BenchmarkCount(b *testing.B) {
	s := "go is simple and fast"
	for i := 0; i < b.N; i++ {
		_ = Count(s)
	}
}
