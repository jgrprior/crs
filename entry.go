package crs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type actionField string

func (af *actionField) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	default:
		return fmt.Errorf("submitAction must be one of 'email' or 'store', found '%v'", s)
	case "email":
		*af = "email"
	case "store":
		*af = "store"
	case "":
		*af = ""
	}

	return nil
}

// An Entry encapsulates submitted entrant and contact permission data, plus
// arbitrary key/value pairs, which are most often mapped to web form fields.
type Entry struct {
	CampaignName       string      `json:"campaignName"`
	CampaignVersion    string      `json:"campaignVersion"`
	SessionFingerprint string      `json:"sessionFingerprint,omitempty"`
	SubmitAction       actionField `json:"submitAction"`
	Entrant            Entrant     `json:"entrant"`
	Permissions        Perms       `json:"permisions"`
	Form               []EntryItem `json:"form,omitempty"`
	Tags               []EntryItem `json:"tags,omitempty"`
	PublicID           string      `json:"entryId,omitempty"`
}

// An Entrant contains name and contact information for a single individual who,
// either through a web for or other proxy to a web service, has submitted their
// details as part of a survey, competition os similar data catpure initiative.
type Entrant struct {
	Title        string `json:"title"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	BirthDate    string `json:"birthDate,omitempty"`
	PhoneNumber  string `json:"phoneNumber,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
}

// A Perms object holds boolean flags that indicate an entrant's contact
// permission preferences at time of entry data capture.
type Perms struct {
	OptInEmail bool `json:"optInEmail,omitempty"`
	OptInPhone bool `json:"optInPhone,omitempty"`
	OptInSms   bool `json:"optInESms,omitempty"`
	OptInPost  bool `json:"optInPost,omitempty"`
}

// An EntryItem holds key/value pairs that are not entrant contact info or
// permissions. Questionnaire or survey questions and answers are examples of
// items submitted as EntryItem.
type EntryItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewEntry initialises a new campaign Entry given a request body.
func NewEntry(body io.Reader) (Entry, error) {
	decoder := json.NewDecoder(body)
	var entry Entry
	err := decoder.Decode(&entry)
	if err == nil {
		entry.PublicID = entryHash(&entry)
	}
	return entry, err
}

func entryHash(e *Entry) string {
	h := md5.New()
	io.WriteString(h, e.Entrant.EmailAddress)
	io.WriteString(h, time.Now().String())
	return hex.EncodeToString(h.Sum(nil))
}

// Valid validates Entry fields and fields of all nested structs. If the Entry
// or any of it's nested structs is not valid, false and a slice of error
// messages is returned.
//
// Entry validation accumulates all error messages from invalid nested structs
// into a single slice of strings.
func (e *Entry) Valid() (bool, []string) {
	valid := true
	msgs := make([]string, 0)

	if e.CampaignName == "" {
		valid = false
		msgs = append(msgs, "campaignName is required")
	}
	if e.CampaignVersion == "" {
		valid = false
		msgs = append(msgs, "campaignVersion is required")
	}

	if res, msg := e.Entrant.Valid(); res == false {
		valid = false
		msgs = append(msgs, msg...)
	}

	if res, msg := e.Permissions.Valid(); res == false {
		valid = false
		msgs = append(msgs, msg...)
	}

	for _, itm := range e.Form {
		if res, msg := itm.Valid(); res == false {
			valid = false
			msgs = append(msgs, msg...)
		}
	}

	for _, itm := range e.Tags {
		if res, msg := itm.Valid(); res == false {
			valid = false
			msgs = append(msgs, msg...)
		}
	}

	return valid, msgs
}

// Valid validates Entrant fields. If the Entrant struct is not valid, false
// and a slice of error messages is returned, otherwise true and an empty slice
// of strings.
func (e *Entrant) Valid() (bool, []string) {
	valid := true
	msgs := make([]string, 0)

	if e.Title == "" {
		valid = false
		msgs = append(msgs, "entrant.title is a required field")
	}
	if e.FirstName == "" {
		valid = false
		msgs = append(msgs, "entrant.firstName is a required field")
	}
	if e.LastName == "" {
		valid = false
		msgs = append(msgs, "entrant.lastName is a required field")
	}
	if e.EmailAddress == "" {
		valid = false
		msgs = append(msgs, "entrant.emailAddress is a required field")
	}

	return valid, msgs
}

// Valid validates Perms fields. If the Perms struct is not valid, false
// and a slice of error messages is returned, otherwise true and an empty slice
// of strings.
func (p *Perms) Valid() (bool, []string) {
	// TODO: implement
	return true, make([]string, 0)
}

// Valid validates EntryItem fields. If the EntryItem struct is not valid,
// false and a slice of error messages is returned, otherwise true and an empty
// slice of strings.
func (i *EntryItem) Valid() (bool, []string) {
	valid := true
	msgs := make([]string, 0)

	if i.Key == "" {
		valid = false
		msgs = append(msgs, "form/tag.key is a required field")
	}
	if i.Value == "" {
		valid = false
		msgs = append(msgs, "form/tag.value is a required field")
	}
	return valid, msgs
}
