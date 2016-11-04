package capture

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	PayloadContextKey = "JsonPayload"
	DefaultDatabase   = "capture"
	DefaultCollection = "entry"
)

var (
	db         string = os.Getenv("GOCAPTURE_DB")
	collection string = os.Getenv("GOCAPTURE_COLLECTION")
)

func init() {
	if db == "" {
		db = DefaultDatabase
	}
	if collection == "" {
		collection = DefaultCollection
	}
}

// An Entry encapsulates submitted entrant and contact permission data, plus
// arbitrary key/value pairs, which are most often mapped to web form fields.
type Entry struct {
	CampaignName       string      `json:"campaignName"`
	CampaignVersion    string      `json:"campaignVersion"`
	SessionFingerprint string      `json:"sessionFingerprint,omitempty"`
	SubmitAction       string      `json:"submitAction"`
	Entrant            Entrant     `json:"entrant"`
	Permissions        Perms       `json:"permisions"`
	Form               []EntryItem `json:"form,omitempty"`
	Tags               []EntryItem `json:"tags,omitempty"`
	PublicId           string      `json:"entryId,omitempty"`
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

func NewEntry(body io.ReadCloser) (Entry, error) {
	decoder := json.NewDecoder(body)
	var entry Entry
	err := decoder.Decode(&entry)
	if err == nil {
		entry.PublicId = entryHash(&entry)
	}
	return entry, err
}

func entryHash(e *Entry) string {
	h := md5.New()
	io.WriteString(h, e.Entrant.EmailAddress)
	io.WriteString(h, time.Now().String())
	return hex.EncodeToString(h.Sum(nil))
}

// Save commits an Entry to a MongoDb database, pre configured by `session`, and
// returns the error, if any, from the mgo.Session call to Insert.
func (e *Entry) Save(s *mgo.Session) error {
	defer s.Close()
	c := s.DB(db).C(collection)
	return c.Insert(e)
}

// IsValid validates Entry fields and fields of all nested structs. If the Entry
// or any of it's nested structs is not valid, false and a slice of error
// messages is returned.
//
// Entry validation accumulates all error messages from invalid nested structs
// into a single slice of strings.
func (e *Entry) Valid() (bool, []string) {
	var (
		valid bool
		msgs  []string
	)

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

// IsValid validates Entrant fields. If the Entrant struct is not valid, false
// and a slice of error messages is returned, otherwise true and an empty slice
// of strings.
func (e *Entrant) Valid() (bool, []string) {
	// TODO: implement
	return true, make([]string, 0)
}

// IsValid validates Perms fields. If the Perms struct is not valid, false
// and a slice of error messages is returned, otherwise true and an empty slice
// of strings.
func (p *Perms) Valid() (bool, []string) {
	// TODO: implement
	return true, make([]string, 0)
}

// IsValid validates EntryItem fields. If the EntryItem struct is not valid,
// false and a slice of error messages is returned, otherwise true and an empty
// slice of strings.
func (i *EntryItem) Valid() (bool, []string) {
	// TODO: implement
	return true, make([]string, 0)
}

// TODO: Move to utils
func JsonEncode(s interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(s)
	return buffer.String(), err
}
