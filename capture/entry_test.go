package capture

import "testing"

func TestValidEntrant(t *testing.T) {

	type test struct {
		input *Entrant
		want  bool
	}

	var tests = []test{
		test{&Entrant{Title: "Mr", FirstName: "John", LastName: "Smith"}, true},
		test{&Entrant{Title: "Mr", FirstName: "Joan", LastName: "Smith"}, true},
	}

	for _, e := range tests {
		if got, _ := e.input.Valid(); got != e.want {
			s, _ := JsonEncode(e.input)
			t.Errorf("Entrant%s did not validate as expected", s)
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
			s, _ := JsonEncode(tst.input)
			t.Errorf("Perms%s did not validate as expected", s)
		}
	}
}
