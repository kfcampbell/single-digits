package parser

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	var tests = []struct {
		input    string
		expected time.Duration
	}{
		{`
		at!) Verizon © 11:45 AM @ 73% a)
a
Congrats!
You solved a Mini puzzle in 1:42.
Have you played our new matching game?
It’s mesmerizing.
TRY TILES
		`, time.Minute*1 + time.Second*42},
		{`
		al! Verizon © 7:01PM a74% 7)
Congrats!
You solved a Mini puzzle in 42 seconds.
Have you played our new matching game?
It’s mesmerizing.
ata a)
		`, time.Second * 42},
		{`
0:57
		`, time.Second * 57},
		{`
Congratulations!
You solved the puzzle in 0:36!
		`, time.Second * 36},
		{`
T-Mobile Wi-Fi + 8:09 PM 69% =)
: f+ .
Congrats!
You finished a Mini puzzle in
0:23
		`, time.Second * 23},
		{`
			C ions!\nongratulations!\nYou solved a Mini puzzle in 26 seconds. :\nGame for something new?\nMake as many words as you can\nwith 7 letters.\nlI O r4
		`, time.Second * 26},
	}

	for _, test := range tests {
		actual, err := GetScoreFromText(test.input)
		if err != nil {
			t.Errorf("Parse(%q) error: %s", test.input, err)
		}
		if actual != test.expected {
			t.Errorf("Parse(%q) = %d, expected %d", test.input, actual, test.expected)
		}
	}
}
