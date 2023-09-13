package jsoner_test

import (
	"testing"

	"github.com/gomig/jsoner"
)

func TestJsoner(t *testing.T) {
	type Book struct {
		Title string
		ISBN  string `json:"isbn"`
	}
	type Author struct {
		Name      string `json:"name"`
		Family    string `json:"family"`
		Age       int    `json:"age,string,omitempty"`
		IsMariage bool   `json:",string"`
		Books     []Book `json:"author_books"`
		Skills    []string
		Address   map[string]any
		// ignored fields
		PrivateField string `json:"-"`
		ignored      string
	}

	john := Author{
		Name:      "John",
		Family:    "Doe",
		Age:       0,
		IsMariage: false,
		Books: []Book{
			{Title: "Basics Of C", ISBN: "12345"},
			{Title: "Golang", ISBN: "88888"},
		},
		Skills: []string{"Web dev", "System programming", "IOT"},
		Address: map[string]any{
			"state": map[string]string{
				"country": "USA",
				"county":  "NY",
			},
			"city":   "NY city",
			"street": "ST. 23",
			"no":     13,
		},
		PrivateField: "Some private information",
		ignored:      "i'm ignored",
	}
	options := map[string]jsoner.JsonerOption{
		".":             {Ignore: []string{"family"}},
		"Address.state": {Ignore: []string{"country"}},
		"author_books":  {Only: []string{"Title"}},
	}
	bytes, _ := jsoner.Jsoner(&john, options)
	excpt := `{"Address":{"city":"NY city","no":13,"state":{"county":"NY"},"street":"ST. 23"},"IsMariage":"false","Skills":["Web dev","System programming","IOT"],"author_books":[{"Title":"Basics Of C"},{"Title":"Golang"}],"name":"John"}`
	if excpt != string(bytes) {
		t.Log(`want: ` + excpt)
		t.Log(`get: ` + string(bytes))
		t.Fatal("not excepted result!")
	}
}
