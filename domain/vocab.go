package domain

import (
	"fmt"
	"sort"
	"strings"
)

type Vocab struct {
	ID     int
	UserID int
}

func (v *Vocab) String() string {
	return fmt.Sprintf("Vocab (ID: %v; UserID: %v)", v.ID, v.UserID)
}

type VocabEntry struct {
	ID            int
	Text          string
	Transcription string
	Translations  []*Translation
}

func (e *VocabEntry) String() string {
	builder := new(strings.Builder)
	builder.WriteString(fmt.Sprintf("VocabEntry (ID: %v; Text: %s; Transcription: %s; Translations: [",
		e.ID, e.Text, e.Transcription))
	for i, t := range e.Translations {
		if i != 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(t.String())
	}
	builder.WriteString("])")
	return builder.String()
}

func (e *VocabEntry) ShortDesc() string {
	var translation string
	for _, t := range e.Translations {
		if t.Position == 0 {
			translation = t.Text
			break
		}
	}
	return fmt.Sprintf("[%s]\n%s", e.Transcription, translation)
}

func (e *VocabEntry) FullDesc() string {
	builder := new(strings.Builder)
	builder.WriteString(e.ShortDesc())
	sort.Slice(e.Translations, func(i, j int) bool {
		return e.Translations[i].Position < e.Translations[j].Position
	})
	var lastClass string
	for _, t := range e.Translations {
		if t.Position == 0 {
			continue
		}
		if t.Class != lastClass {
			builder.WriteString("\n\n" + t.Class + ": " + t.Text)
			lastClass = t.Class
		} else {
			builder.WriteString(", " + t.Text)
		}
	}
	return builder.String()
}

type Translation struct {
	ID       int
	Text     string
	Class    string
	Position int
}

func (t *Translation) String() string {
	return fmt.Sprintf("Translation (ID: %v, Text: %s, Class: %s, Position: %v", t.ID, t.Text, t.Class, t.Position)
}
