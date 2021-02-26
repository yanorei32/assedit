package assutil

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type status int

const (
	BEFORE_EVENTS status = iota
	FIND_EVENT_FORMAT
	EVENTS
	AFTER_EVENTS

	CAPTION2ASS_FORMAT = "Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\r\n"
)

func Parse(f *bufio.Reader) (*Ass, error) {
	ass := Ass{}
	s := BEFORE_EVENTS

	for {
		l, err := f.ReadString('\n')

		if err == io.EOF {
			if s != EVENTS && s != AFTER_EVENTS {
				return nil, errors.New("Unexpected EOF")
			}

			break
		}

		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(l, "[") {
			switch s {
			case BEFORE_EVENTS:
				if l == "[Events]\r\n" {
					s = FIND_EVENT_FORMAT
				}
			case EVENTS:
				s = AFTER_EVENTS

			}
		}

		switch s {
		case BEFORE_EVENTS:
			ass.OtherSectionsHead += l

		case FIND_EVENT_FORMAT:
			if strings.HasPrefix(l, "Format:") {
				if l != CAPTION2ASS_FORMAT {
					return nil, errors.New("Unexpected ASS format")
				}

				s = EVENTS
			}
			ass.OtherSectionsHead += l

		case EVENTS:
			if strings.HasPrefix(l, "Dialogue:") {
				e, err := NewEvent(l)

				if err != nil {
					return nil, err
				}

				ass.EventSection = append(ass.EventSection, *e)
			}

		case AFTER_EVENTS:
			ass.OtherSectionsTail += l

		}
	}

	return &ass, nil
}
