package combatlog

import (
	"testing"
	"time"
)

// Groups of 5 every 3
var testCombatLog = CombatLog{

	// 0..4                                              10..14
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 10 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 11 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 12 }, },
	// 3..5                                              13..17
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 13 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 14 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 15 }, },
	// 6..6                                              16..20
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 17 }, },
	// 7..7                                              19..23
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 20 }, },
	// 8..8                                              22..26
	// 8..8                                              25..29
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 25 }, },
	//                                                   28..32
	//                                                   31..35
	//                                                   34..38
	//                                                   37..41
	//                                                   40..44
	//                                                   43..47
	//                                                   46..50
	//                                                   49..53
	// 9..10                                             52..56
	// 9..13                                            55..59
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 55 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 56 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 57 }, },
	// 12..17                                            58..62
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 58 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 1, Second: 59 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 2, Second: 00 }, },
	// 15..17                                            01..05
	Event{ Time: &time.Time{ Year: 2011, Minute: 2, Second: 01 }, },
	Event{ Time: &time.Time{ Year: 2011, Minute: 2, Second: 02 }, },
}

var windowTests = []struct {
	LastStart, Start, LastEnd, End int
}{
	{0, 0, 0, 5},
	{0, 3, 5, 6},
	{3, 6, 6, 7},
	{6, 7, 7, 8},
	{7, 8, 8, 9},
	{8, 8, 9, 9},
	{8, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 9},
	{9, 9, 9, 11},
	{9, 9, 11, 14},
	{9, 12, 14, 17},
	{12, 15, 17, 17},
	{15, 17, 17, 17},
}

func TestWindow(t *testing.T) {
	idx := 0
	window := func(ls, s, le, e int) {
		if idx >= len(windowTests) {
			t.Errorf("extra window(%d, %d, %d, %d)", ls, s, le, e)
			return
		}
		test := windowTests[idx]
		if ls != test.LastStart || s != test.Start || false {
			//le != test.LastEnd || e != test.End {
			t.Errorf("%d. got window(%d, %d, %d, %d), want window(%d, %d, %d, %d)",
				idx, ls, s, le, e, test.LastStart, test.Start, test.LastEnd, test.End)
		} else {
			t.Logf("%d. matched window(%d, %d, %d, %d)", idx, ls, s, le, e)
		}

		idx++
	}
	testCombatLog.RollingWindow(5, 3, window)
	for i := idx; i < len(windowTests); i++ {
		test := windowTests[i]
		t.Errorf("missing window(%d, %d, %d, %d)",
			test.LastStart, test.Start, test.LastEnd, test.End)
	}
}
