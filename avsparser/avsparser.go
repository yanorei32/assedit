package avsparser

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func Parse(f *bufio.Reader) (*Avs, error) {
	avs := Avs{}

	regex := regexp.MustCompile(`Trim\((?P<start>\d+),(?P<end>\d+)\)`)

	for {
		l, err := f.ReadString('\n')

		if err == io.EOF {
			return nil, errors.New("Unexpected EOF")
		}

		if err != nil {
			return nil, err
		}

		if !strings.HasPrefix(l, "Trim") {
			continue
		}

		matches := regex.FindAllStringSubmatch(l, -1)
		for _, m := range matches {
			trim := Trim{}

			for i, v := range m {
				vUint, _ := strconv.ParseUint(v, 10, 64)

				switch regex.SubexpNames()[i] {
				case "start":
					trim.StartFrame = vUint
				case "end":
					trim.EndFrame = vUint
				}

			}

			avs.Trims = append(avs.Trims, trim)
		}

		return &avs, nil
	}
}
