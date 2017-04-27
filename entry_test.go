package crs

import (
	"bytes"
	"encoding/json"
	"testing"
)

func jsonEncode(s interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(s)
	return buffer.String(), err
}

func TestValidEntrant(t *testing.T) {

	type test struct {
		input *Entrant
		want  bool
	}

	var tests = []test{
		test{&Entrant{Title: "Mr", FirstName: "John", LastName: "Smith", EmailAddress: "foo@bar.com"}, true},
		test{&Entrant{Title: "Mr", FirstName: "Joan", LastName: "Smith", EmailAddress: "foo@bar.com"}, true},
	}

	for _, e := range tests {
		if got, _ := e.input.Valid(); got != e.want {
			s, _ := jsonEncode(e.input)
			t.Errorf("Entrant%s did not validate as expected", s)
		}
	}
}

func TestValidEntry(t *testing.T) {

	type test struct {
		input *Entry
		want  bool
	}

	var tests = []test{
		test{&Entry{
			CampaignName:    "Foo",
			CampaignVersion: "0.0.1",
			Entrant: Entrant{
				Title:        "Mr",
				FirstName:    "John",
				LastName:     "Smith",
				EmailAddress: "foo@bar.com",
			},
		},
			true},
	}

	for _, e := range tests {
		if got, errs := e.input.Valid(); got != e.want {
			//s, _ := jsonEncode(e.input)
			t.Errorf("Entry %s did not validate as expected", errs)
		}
	}
}

func TestValidPerms(t *testing.T) {

	type test struct {
		input *Perms
		want  bool
	}

	var tests = []test{
		test{&Perms{OptInEmail: true, OptInPhone: true, OptInSms: true, OptInPost: true}, true},
	}

	for _, tst := range tests {
		if got, _ := tst.input.Valid(); got != tst.want {
			s, _ := jsonEncode(tst.input)
			t.Errorf("Perms%s did not validate as expected", s)
		}
	}
}
