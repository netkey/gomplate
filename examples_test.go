package gomplate

import (
	"bytes"
	"fmt"
)

// Demonstrates how to render a simple template which uses gomplate functions but no datasources.
func ExampleRenderTemplate_noDatasources() {
	inString := `{{ slice "banana" "cheese" "apple" "donut" | coll.Sort }}`

	in := bytes.NewBufferString(inString)
	out := &bytes.Buffer{}

	err := RenderTemplate(in, out, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(out.String())

	// Output: [apple banana cheese donut]
}

// Options can be provided to customize gomplate's behaviour.
// This example shows how to use custom delimiters.
func ExampleRenderTemplate_withOptions() {
	inString := `{% strings.ToUpper "hello world" }{{ unrendered }}`

	in := bytes.NewBufferString(inString)
	out := &bytes.Buffer{}

	opts := &Options{
		LDelim: "{%",
		RDelim: "}",
	}

	err := RenderTemplate(in, out, opts)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(out.String())

	// Output: HELLO WORLD{{ unrendered }}
}

// Datasources can also be provided.
func ExampleRenderTemplate_datasources() {
	inString := `{{ range (ds "myData").values -}}
{{ . }}
{{ end }}`

	// myData := &data.Source{
	// 	Alias: "myData",
	// }

	in := bytes.NewBufferString(inString)
	out := &bytes.Buffer{}

	opts := &Options{
		// DataSources: []DataSource{myData},
	}

	err := RenderTemplate(in, out, opts)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(out.String())

	// Output: HELLO WORLD{{ unrendered }}
}
