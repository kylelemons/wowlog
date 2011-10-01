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
	ReadBufferSize  = 4*1024*1024 // 4MB
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
	//return Read(file)
	cl, err := Read(file)
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(cl); err != nil {
		return nil, err
	}
	fmt.Println("%s: encodes to %d bytes", buf.Len())
	return cl, err
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

type eventField struct {
	index        []int
	reflect.Kind
}

type eventFactory struct {
	fields   []eventField
	min, max int
	reflect.Type
}

func (e eventFactory) String() string {
	return fmt.Sprintf("%s%v", e.Type, e.fields)
}

func (e eventFactory) create(csv []byte) (interface{}, os.Error) {
	start := 0
	parsed := 0

	val := reflect.New(e.Type)
	elem := val.Elem()

	for i, field := range e.fields {
		if start >= len(csv) {
			break
		}
		comma := nextField(csv[start:])
		fstr := string(csv[start:][:comma])

		fval := elem.FieldByIndex(field.index)
		if err := parseField(fval, field.Kind, fstr); err != nil {
			return nil, fmt.Errorf("combatlog: failed to parse %q as %s for %s[%d]: %s",
				fstr, field.Kind, e.Type, i, err)
		}

		parsed++
		start += comma + 1
	}

	if parsed < e.min || parsed > e.max {
		return nil, fmt.Errorf("combatlog: event has %d fields, it should have %d-%d:\n%s\n%s",
			parsed, e.min, e.max, e, string(csv))
	}

	return val.Interface(), nil
}

func parseField(val reflect.Value, kind reflect.Kind, fstr string) (err os.Error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
		}
	}()

	switch kind {
	case reflect.Int64, reflect.Int32:
		i, err := strconv.Btoi64(fstr, 0)
		if err != nil {
			return err
		}
		val.SetInt(i)
	case reflect.Uint64, reflect.Uint32:
		u, err := strconv.Btoui64(fstr, 0)
		if err != nil {
			return err
		}
		val.SetUint(u)
	case reflect.String:
		if fstr[0] == '"' {
			if fstr, err = strconv.Unquote(fstr); err != nil {
				return err
			}
		}
		val.SetString(fstr)
	case reflect.Bool:
		switch fstr {
		case "nil":
			val.SetBool(false)
		case "1":
			val.SetBool(true)
		default:
			return fmt.Errorf("invalid bool: %q", fstr)
		}
	default:
		return fmt.Errorf("unknown kind %s", kind)
	}
	return nil
}

func compile(empty interface{}) (comp eventFactory) {
	gob.Register(empty)

	eventType := reflect.TypeOf(empty)
	for eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	comp.Type = eventType

	optional := false

	var next func([]int, reflect.Type)
	next = func(curr []int, t reflect.Type) {
		for i, n := 0, t.NumField(); i < n; i++ {
			// Get the recursive index
			idx := make([]int, len(curr)+1)
			idx[copy(idx, curr)] = i

			ft := t.Field(i)
			if ft.Tag.Get("combatlog") == "optional" {
				optional = true
			}

			if fk := ft.Type.Kind(); fk == reflect.Struct {
				next(idx, ft.Type)
			} else {
				if !optional {
					comp.min++
				}
				comp.max++
				comp.fields = append(comp.fields, eventField{idx,fk})
			}
		}
	}

	next(nil, eventType)
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
				case '"' :
					break closeQuote
				}
			}
		}
	}
	return len(csv)
}
