package assutil

import (
	"errors"
	"fmt"
	"time"
)

type Event struct {
	TextPart1 string
	StartTime time.Duration
	TextPart2 string
	EndTime   time.Duration
	TextPart3 string
}

func asstime2duration(s string) (time.Duration, error) {
	return time.ParseDuration(
		fmt.Sprintf("%sh%sm%ss", s[0:1], s[2:4], s[5:]),
	)
}

func duration2asstime(d time.Duration) string {
	d = d.Round(time.Millisecond * 10)

	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	cs := d / (time.Millisecond * 10)

	return fmt.Sprintf("%01d:%02d:%02d.%02d", h, m, s, cs)
}

func NewEvent(l string) (*Event, error) {
	e := Event{}

	e.TextPart1 = l[:12]

	start, err := asstime2duration(l[12:22])
	e.StartTime = start
	if err != nil {
		return nil, errors.New("Unexpected ASS format")
	}

	e.TextPart2 = l[22:23]

	end, err := asstime2duration(l[23:33])
	e.EndTime = end
	if err != nil {
		return nil, errors.New("Unexpected ASS format")
	}

	e.TextPart3 = l[33:]

	return &e, nil
}

func (e Event) Format() (s string) {
	s += e.TextPart1
	s += duration2asstime(e.StartTime)
	s += e.TextPart2
	s += duration2asstime(e.EndTime)
	s += e.TextPart3

	return
}
