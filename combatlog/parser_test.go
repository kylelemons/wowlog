package combatlog

import (
	"bytes"
	"io"
	"path/filepath"
	"reflect"
	"time"
	"testing"
)

var realLogs = glob("CombatLog.*.txt")

func glob(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic("glob: " + err.String())
	}
	if len(files) == 0 {
		panic("glob: no files: " + pattern)
	}
	return files
}

func TestRealCombatLog(t *testing.T) {
	return
	for _, file := range realLogs {
		cl, err := ReadFile(file)
		if err != nil {
			t.Errorf("load(%q): %s", file, err)
		}
		_ = cl
	}
}

var decodeTests = []struct {
	Desc   string
	Lines  string
	Events CombatLog
}{
	{
		Desc: "re-use timestamps",
		Lines: `
9/25 19:03:22.951  SPELL_DAMAGE,0xF130966900007981,"Knight of the Ebon Blade",0xa18,0x0,0xF15079A30069A7D9,"Pustulent Horror",0xa48,0x0,66019,"Death Coil",0x20,5087,-1,32,0,0,0,nil,nil,nil
9/25 19:03:23.045  SWING_DAMAGE,0xF15096640000699C,"Argent Warhorse",0xa18,0x0,0xF15079A30069A7D9,"Pustulent Horror",0xa48,0x0,155,-1,1,0,0,0,nil,nil,nil
`,
		Events: CombatLog{
			Event{
				Time: time.Time{
					Year: 0, Month: 9, Day: 25,
					Hour: 19, Minute: 3, Second: 22, Nanosecond: 951000000,
					ZoneOffset: 0, Zone: "",
				},
				Name: "SPELL_DAMAGE",
				Data: SpellDamage{
					Common: Common{
						Source: Unit{
							Name: "Knight of the Ebon Blade",
							ID:   0xf130966900007981, Flags: 0xa18, Null: 0,
						},
						Dest: Unit{
							Name: "Pustulent Horror",
							ID:   0xf15079a30069a7d9, Flags: 0xa48, Null: 0,
						},
					},
					Spell: Spell{
						Name: "Death Coil",
						ID:   0x101e3, School: 32,
					},
					Damage: Damage{
						Amount: 5087, School: 32,
						Resisted: 0, Blocked: 0, Absorbed: 0,
						Critical: false, Glancing: false, Crushing: false,
						Unknown: -1,
					},
				},
			},
			Event{
				Time: time.Time{
					Year: 0, Month: 9, Day: 25,
					Hour: 19, Minute: 3, Second: 23, Nanosecond: 45000000,
					ZoneOffset: 0, Zone: "",
				},
				Name: "SWING_DAMAGE",
				Data: SwingDamage{
					Common: Common{
						Source: Unit{
							Name: "Argent Warhorse",
							ID:   0xf15096640000699c, Flags: 0xa18,
						},
						Dest: Unit{
							Name: "Pustulent Horror",
							ID:   0xf15079a30069a7d9, Flags: 0xa48,
						},
					},
					Damage: Damage{
						Amount: 155, School: 1,
						Resisted: 0, Blocked: 0, Absorbed: 0,
						Critical: false, Glancing: false, Crushing: false,
						Unknown: -1,
					},
				},
			},
		},
	},
	{
		Desc: "don't overwrite",
		Lines: `
9/25 19:03:23.045  SWING_DAMAGE,0xF15096640000699C,"Argent Warhorse",0xa18,0x0,0xF15079A30069A7D9,"Pustulent Horror",0xa48,0x0,155,-1,1,0,0,0,nil,nil,nil
9/25 19:03:23.079  SWING_DAMAGE,0xF15079A30069A7D9,"Pustulent Horror",0xa48,0x0,0xF130965D00687234,"Argent Crusader",0xa18,0x0,1820,-1,1,0,0,0,nil,nil,nil
`,
		Events: CombatLog{
			Event{
				Time: time.Time{
					Year: 0, Month: 9, Day: 25,
					Hour: 19, Minute: 3, Second: 23, Nanosecond: 45000000,
					ZoneOffset: 0, Zone: "",
				},
				Name: "SWING_DAMAGE",
				Data: SwingDamage{
					Common: Common{
						Source: Unit{
							Name: "Argent Warhorse",
							ID:   0xf15096640000699c, Flags: 0xa18,
						},
						Dest: Unit{
							Name: "Pustulent Horror",
							ID:   0xf15079a30069a7d9, Flags: 0xa48,
						},
					},
					Damage: Damage{
						Amount: 155, School: 1,
						Resisted: 0, Blocked: 0, Absorbed: 0,
						Critical: false, Glancing: false, Crushing: false,
						Unknown: -1,
					},
				},
			},
			Event{
				Time: time.Time{
					Year: 0, Month: 9, Day: 25,
					Hour: 19, Minute: 3, Second: 23, Nanosecond: 79000000,
					ZoneOffset: 0, Zone: "",
				},
				Name: "SWING_DAMAGE",
				Data: SwingDamage{
					Common: Common{
						Source: Unit{
							Name: "Pustulent Horror",
							ID:   0xf15079a30069a7d9, Flags: 0xa48,
						},
						Dest: Unit{
							Name: "Argent Crusader",
							ID:   0xf130965d00687234, Flags: 0xa18,
						},
					},
					Damage: Damage{
						Amount: 1820, School: 1,
						Resisted: 0, Blocked: 0, Absorbed: 0,
						Critical: false, Glancing: false, Crushing: false,
						Unknown: -1,
					},
				},
			},
		},
	},
}

func TestDecode(t *testing.T) {
	for _, test := range decodeTests {
		cl, err := Read(bytes.NewBufferString(test.Lines))
		if err != nil {
			t.Errorf("%s: error: %s", test.Desc, err)
		}
		if !reflect.DeepEqual(cl, test.Events) {
			t.Logf("%s: got:", test.Desc)
			for _, e := range cl {
				t.Logf("  %#v", e)
			}
			t.Logf("%s: want:", test.Desc)
			for _, e := range test.Events {
				t.Logf("  %#v", e)
			}
			t.Fail()
		}
	}
}

var nextFieldTests = []struct {
	Source string
	Comma  int
}{
	{"a", 1},
	{"a,b", 1},
	{",b", 0},
	{`"b,c",d`, 5},
	{`"b\",c",d`, 7},
}

func TestNextField(t *testing.T) {
	for _, test := range nextFieldTests {
		if got, want := nextField(test.Source), test.Comma; got != want {
			t.Errorf("nextField(%#q) = %d, want %d", test.Source, got, want)
		}
	}
}

var benchDecodeLine = `9/25 19:03:27.752  SPELL_DAMAGE,0xF15079A30069A7D9,"Pustulent Horror",0xa48,0x0,0xF130966900006997,"Knight of the Ebon Blade",0xa18,0x0,28405,"Knockback",0x1,708,-1,1,0,0,0,nil,nil,nil`

func BenchmarkDecode(b *testing.B) {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		line := benchDecodeLine + "\n"
		for i := 0; i < b.N; i++ {
			io.WriteString(w, line)
		}
	}()
	if _, err := Read(r); err != nil {
		panic(err)
	}
}

var benchFile = "CombatLog.Bench.txt"

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFile(benchFile)
	}
}

func BenchmarkLoadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFile(benchFile)
	}
}
