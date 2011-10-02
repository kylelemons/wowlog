package combatlog

import (
	"bufio"
	"bytes"
	"fmt"
	"gob"
	"io"
	"os"
	"reflect"
	"strconv"
	"time"
)

const (
	TimeStampFormat = "1/2 15:04:05.000"
	ReadBufferSize  = 4 * 1024 * 1024 // 4MB
)

type Event struct {
	Time *time.Time
	Name string
	Data interface{}
}

type CombatLog []Event

// filename leaks
func ReadFile(filename string) (CombatLog, os.Error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return Read(file)
	/*
		cl, err := Read(file)
		buf := new(bytes.Buffer)
		if err := gob.NewEncoder(buf).Encode(cl); err != nil {
			return nil, err
		}
		fmt.Println("%s: encodes to %d bytes", buf.Len())
		return cl, err
	*/
}

func start_of_event(rune int) bool { return rune >= '@' }

// r leaks
func Read(r io.Reader) (events CombatLog, err os.Error) {
	lines, err := bufio.NewReaderSize(r, ReadBufferSize)
	if err != nil {
		return nil, err
	}

	for {
		// Read the next line
		line, isPrefix, err := lines.ReadLine()
		if err == os.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Finish the line if it's super long
		if isPrefix {
			/*
				full := make([]byte, 0, len(line))
				full = append(full, line...)
				for isPrefix {
					line, isPrefix, err = lines.ReadLine()
					if err == os.EOF {
						break
					}
					if err != nil {
						return nil, err
					}
					full = append(full, line...)
				}
				line = full
			*/
			return nil, fmt.Errorf("combatlog: long line: %q\n", line)
		}

		// Skip the line if it's blank
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Split the line into pieces
		start := bytes.IndexFunc(line, start_of_event)
		if start < 0 {
			fmt.Printf("combatlog: malformatted line: %q\n", line)
			continue
		}

		stamp := string(bytes.TrimSpace(line[:start]))
		time, err := time.Parse(TimeStampFormat, stamp)

		csv := line[start:]
		comma := bytes.IndexByte(csv, ',')
		if comma < 0 {
			fmt.Printf("combatlog: malformatted line: %q\n", line)
			continue
		}

		name, csv := string(csv[:comma]), csv[comma+1:]
		factory, ok := eventTypes[name]
		if !ok {
			return nil, fmt.Errorf("combatlog: unknown event type %q", name)
		}

		data, err := factory.create(csv)
		if err != nil {
			return nil, err
		}

		event := Event{
			Time: time,
			Name: name,
			Data: data,
		}

		events = append(events, event)
	}

	return events, nil
}

type field interface {
	parse(string) os.Error
	zero()
}

type eventFactory struct {
	fields   []field
	min, max int
	emptyptr interface{}
}

func (e eventFactory) String() string {
	return fmt.Sprintf("%T%v", e.emptyptr, e.fields)
}

func (e eventFactory) create(csv []byte) (event interface{}, err os.Error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()

	start := 0
	parsed := 0

	for i, field := range e.fields {
		if start >= len(csv) {
			break
		}
		comma := nextField(csv[start:])
		fstr := string(csv[start:][:comma])

		if err := field.parse(fstr); err != nil {
			return nil, fmt.Errorf("combatlog: failed to parse %q as %T for %T[%d]: %s",
				fstr, field, e.emptyptr, i, err)
		}

		parsed++
		start += comma + 1
	}

	if parsed < e.min || parsed > e.max {
		return nil, fmt.Errorf("combatlog: event has %d fields, it should have %d-%d:\n%s\n%s",
			parsed, e.min, e.max, e, string(csv))
	}

	for i := parsed; i < e.max; i++ {
		e.fields[i].zero()
	}

	return reflect.ValueOf(e.emptyptr).Elem().Interface(), nil
}

type fieldInt32 struct{ ptr *int32 }

func (f fieldInt32) zero() {
	*f.ptr = 0
}
func (f fieldInt32) parse(fstr string) os.Error {
	i, err := strconv.Btoi64(fstr, 0)
	*f.ptr = int32(i)
	return err
}

type fieldInt64 struct{ ptr *int64 }

func (f fieldInt64) zero() {
	*f.ptr = 0
}
func (f fieldInt64) parse(fstr string) (err os.Error) {
	*f.ptr, err = strconv.Btoi64(fstr, 0)
	return err
}

type fieldUint32 struct{ ptr *uint32 }

func (f fieldUint32) zero() {
	*f.ptr = 0
}
func (f fieldUint32) parse(fstr string) os.Error {
	i, err := strconv.Btoui64(fstr, 0)
	*f.ptr = uint32(i)
	return err
}

type fieldUint64 struct{ ptr *uint64 }

func (f fieldUint64) zero() {
	*f.ptr = 0
}
func (f fieldUint64) parse(fstr string) (err os.Error) {
	*f.ptr, err = strconv.Btoui64(fstr, 0)
	return err
}

type fieldString struct{ ptr *string }

func (f fieldString) zero() {
	*f.ptr = ""
}
func (f fieldString) parse(fstr string) (err os.Error) {
	if fstr[0] == '"' {
		*f.ptr, err = strconv.Unquote(fstr)
	} else {
		*f.ptr = fstr
	}
	return err
}

type fieldBool struct{ ptr *bool }

func (f fieldBool) zero() {
	*f.ptr = false
}
func (f fieldBool) parse(fstr string) (err os.Error) {
	switch fstr {
	case "nil":
		*f.ptr = false
	case "1":
		*f.ptr = true
	default:
		err = fmt.Errorf("invalid bool: %q", fstr)
	}
	return err
}

func compile(empty interface{}) (comp eventFactory) {
	gob.Register(empty)

	val := reflect.ValueOf(empty)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("combatlog: compile: cannot compile a non-pointer-to-struct value")
	}
	comp.emptyptr = val.Interface()
	val = val.Elem()

	optional := false

	var next func([]int, reflect.Value)
	next = func(curr []int, val reflect.Value) {
		typ := val.Type()

		for i, n := 0, val.NumField(); i < n; i++ {
			// Get the recursive index
			idx := make([]int, len(curr)+1)
			idx[copy(idx, curr)] = i

			fval := val.Field(i)
			ftyp := typ.Field(i)
			if ftyp.Tag.Get("combatlog") == "optional" {
				optional = true
			}

			var curField field
			switch ftyp.Type.Kind() {
			case reflect.Struct: next(idx, fval)
			case reflect.Int32:  curField = fieldInt32{fval.Addr().Interface().(*int32)}
			case reflect.Int64:  curField = fieldInt64{fval.Addr().Interface().(*int64)}
			case reflect.Uint32: curField = fieldUint32{fval.Addr().Interface().(*uint32)}
			case reflect.Uint64: curField = fieldUint64{fval.Addr().Interface().(*uint64)}
			case reflect.Bool:   curField = fieldBool{fval.Addr().Interface().(*bool)}
			case reflect.String: curField = fieldString{fval.Addr().Interface().(*string)}
			default:             panic("cannot compile field of type "+ftyp.Type.Kind().String())
			}
			if curField != nil {
				if !optional {
					comp.min++
				}
				comp.max++
				comp.fields = append(comp.fields, curField)
			}
		}
	}

	next(nil, val)
	return comp
}

func nextField(csv []byte) int {
	for i, n := 0, len(csv); i < n; i++ {
		switch csv[i] {
		case ',':
			return i
		case '\\':
			i++
		case '"':
		closeQuote:
			for i++; i < n; i++ {
				switch csv[i] {
				case '\\':
					i++
				case '"':
					break closeQuote
				}
			}
		}
	}
	return len(csv)
}
