# Jsoner

Utility function for marshal json with pattern.

You can filter nested field with path using `.` separator (e.g. `user.address` will apply pattern to `{"user" : { "address" : {...} } }`.

**NOTE:** use "." as path pattern for root object.

**NOTE:** you can use "JsonerIndent" method to generate json with indent.

```go
import "github.com/gomig/jsoner"
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
    // ignored fields
    PrivateField string `json:"-"`
    ignored      string
}

john := Author{
    Name: "John",
    Family: "Doe",
    Age: 0,
    IsMariage: false,
    Books: []Book{
        {Title: "Basics Of C", ISBN: "12345"},
        {Title: "Golang", ISBN: "88888"},
    },
    Skills: []string{"Web dev", "System programming", "IOT"},
    PrivateField: "Some private information",
    ignored: "i'm ignored"
}
options := map[string]jsoner.JsonerOption{
    ".":            {Ignore: []string{"family"}}, // ignore family field from root struct
    "author_books": {Only: []string{"title"}},    // only marshal title field of author books
}
bytes, _ := jsoner.Jsoner(john, options)
fmt.Println(string(bytes))

/*
{
   "author_books": [
      {
         "Title": "Basics Of C"
      },
      {
         "Title": "Golang"
      }
   ],
   "IsMariage": "false",
   "name": "John",
   "Skills": [
      "Web dev",
      "System programming",
      "IOT"
   ]
}
*/
```
