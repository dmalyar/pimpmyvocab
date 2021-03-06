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
	return fmt.Sprintf("ID: %v; UserID: %v", v.ID, v.UserID)
}

type VocabEntry struct {
	ID              int
	Text            string
	Transcription   string
	MainTranslation string
	Translations    []*Translation
}

func (e *VocabEntry) String() string {
	builder := new(strings.Builder)
	builder.WriteString(fmt.Sprintf("ID: %v; Text: %s; Transcription: %s; MainTranslation: %s; Translations: [",
		e.ID, e.Text, e.Transcription, e.MainTranslation))
	for i, t := range e.Translations {
		if i != 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(t.String())
	}
	builder.WriteString("]")
	return builder.String()
}

func (e *VocabEntry) ShortDesc() string {
	if e.Transcription != "" {
		return fmt.Sprintf("[%s]\n%s", e.Transcription, e.MainTranslation)
	} else {
		return e.MainTranslation
	}
}

func (e *VocabEntry) FullDesc(printEntryText bool) string {
	builder := new(strings.Builder)
	if printEntryText {
		builder.WriteString(fmt.Sprintln(e.Text))
	}
	if e.Transcription != "" {
		builder.WriteString(fmt.Sprintf("[%s]", e.Transcription))
	}
	sort.Slice(e.Translations, func(i, j int) bool {
		return e.Translations[i].Position < e.Translations[j].Position
	})
	var lastClass string
	for _, t := range e.Translations {
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
	return fmt.Sprintf("ID: %v, Text: %s, Class: %s, Position: %v", t.ID, t.Text, t.Class, t.Position)
}
