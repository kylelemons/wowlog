package combatlog

import (
	"testing"
	"path/filepath"
	"io"
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
	for _, file := range realLogs {
		cl, err := Load(file)
		if err != nil {
			t.Errorf("load(%q): %s", file, err)
		}
		_ = cl
	}
}

var countFieldTests = []struct{
	Source string
	Count  int
}{
	{"a", 1},
	{"a,b", 2},
	{"a,,b", 3},
	{`a,"b,c",d`, 3},
}

func TestCountFields(t *testing.T) {
	for _, test := range countFieldTests {
		if got, want := countFields([]byte(test.Source)), test.Count; got != want {
			t.Errorf("countFields(%#q) = %d, want %d", test.Source, got, want)
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
	if _, err := ReadFrom(r); err != nil {
		panic(err)
	}
}
