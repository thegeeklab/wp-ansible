// Command docs-gen generates docs/data/data.yaml from the plugin's flag
// definitions and the long descriptions attached as doc comments above each
// flag literal.
//
// Invoke via:  go generate ./plugin/...
package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/thegeeklab/wp-ansible/plugin"
	plugin_docs "github.com/thegeeklab/wp-plugin-go/v6/docs"
	plugin_template "github.com/thegeeklab/wp-plugin-go/v6/template"
)

func main() {
	outputFile := flag.String("output", "", "Output file path")
	sourceFile := flag.String("source", "plugin.go", "Plugin source file to scan for long descriptions")

	flag.Parse()

	if *outputFile == "" {
		log.Fatal("no output file specified")
	}

	p := plugin.New(nil)
	templateData := plugin_docs.GetTemplateDataWithSource(p.App, *sourceFile)
	longDescs := plugin_docs.LongDescriptionsFor(*sourceFile)

	funcs := plugin_template.LoadFuncMap()
	funcs["longDesc"] = longDescFunc(longDescs)
	funcs["yamlDesc"] = yamlDesc

	docTemplate, err := template.New("docs").Funcs(funcs).Parse(docsTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if err := docTemplate.Execute(&buf, templateData); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(*outputFile, buf.Bytes(), 0o600); err != nil {
		log.Fatal(err)
	}
}

const docsTemplate = `---
{{- if .GlobalArgs }}
properties:
{{- range $v := .GlobalArgs }}
  - name: {{ $v.Name }}
    {{- with longDesc $v }}
    description: |
      {{ . | ToSentence | yamlDesc }}
    {{- end }}
    {{- with $v.Type }}
    type: {{ . }}
    {{- end }}
    {{- with $v.Default }}
    defaultValue: {{ . }}
    {{- end }}
    required: {{ default false $v.Required }}
{{ end -}}
{{ end -}}
`

func yamlDesc(s string) string {
	if !strings.Contains(s, "\n") {
		return s
	}

	indented := strings.ReplaceAll(s, "\n", "\n      ")

	return "\n      " + indented
}

func longDescFunc(longDescs map[string]string) func(*plugin_docs.PluginArg) string {
	return func(arg *plugin_docs.PluginArg) string {
		if d, ok := longDescs[arg.Name]; ok && d != "" {
			return d
		}

		return arg.Description
	}
}
