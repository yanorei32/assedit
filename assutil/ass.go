package assutil

type Ass struct {
	OtherSectionsHead string
	EventSection      []Event
	OtherSectionsTail string
}

func (a Ass) Format() (s string) {
	s += a.OtherSectionsHead

	for _, e := range a.EventSection {
		s += e.Format()
	}

	s += a.OtherSectionsTail

	return
}
