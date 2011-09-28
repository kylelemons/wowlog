package combatlog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

const TimeStampFormat = "1/2 15:04:05.000"

type Event struct {
	Time *time.Time
	Name string
	Data interface{}
}

type CombatLog []*Event

func Load(filename string) (CombatLog, os.Error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ReadFrom(bytes.NewBuffer(contents))
}

func start_of_event(rune int) bool { return rune >= '@' }

func ReadFrom(r io.Reader) (events CombatLog, err os.Error) {
	lines := bufio.NewReader(r)

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
			fmt.Printf("combatlog: long line: %q\n", line)
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

		event := &Event{
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
	if got := countFields(csv); got < e.min || got > e.max {
		return nil, fmt.Errorf("combatlog: event has %d fields, it should have %d-%d:\n%s\n%s",
			got, e.min, e.max, e, csv)
	}
	return nil, nil
}

func compile(empty interface{}) (comp eventFactory) {
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

func countFields(csv []byte) (fields int) {
	for i, n := 0, len(csv); i < n; i++ {
		switch csv[i] {
		case ',':
			fields++
		case '"':
			for i++; i < n; i++ {
				if csv[i] == '"' {
					break
				}
			}
		}
	}
	return fields+1
}
