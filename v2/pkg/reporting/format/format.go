package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/projectdiscovery/nuclei/v2/pkg/output"
	"github.com/projectdiscovery/nuclei/v2/pkg/types"
)

// Summary returns a formatted built one line summary of the event
func Summary(event *output.ResultEvent) string {
	template := GetMatchedTemplate(event)

	builder := &strings.Builder{}
	builder.WriteString("[")
	builder.WriteString(template)
	builder.WriteString("] [")
	builder.WriteString(types.ToString(event.Info["severity"]))
	builder.WriteString("] ")
	builder.WriteString(types.ToString(event.Info["name"]))
	builder.WriteString(" found on ")
	builder.WriteString(event.Host)
	data := builder.String()
	return data
}

// MarkdownDescription formats a short description of the generated
// event by the nuclei scanner in Markdown format.
func MarkdownDescription(event *output.ResultEvent) string {
	template := GetMatchedTemplate(event)
	builder := &bytes.Buffer{}
	builder.WriteString("**Details**: **")
	builder.WriteString(template)
	builder.WriteString("** ")
	builder.WriteString(" matched at ")
	builder.WriteString(event.Host)
	builder.WriteString("\n\n**Protocol**: ")
	builder.WriteString(strings.ToUpper(event.Type))
	builder.WriteString("\n\n**Full URL**: ")
	builder.WriteString(event.Matched)
	builder.WriteString("\n\n**Timestamp**: ")
	builder.WriteString(event.Timestamp.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	builder.WriteString("\n\n**Template Information**\n\n| Key | Value |\n|---|---|\n")
	for k, v := range event.Info {
		builder.WriteString(fmt.Sprintf("| %s | %s |\n", k, v))
	}
	if event.Request != "" {
		builder.WriteString("\n**Request**\n\n```http\n")
		builder.WriteString(event.Request)
		builder.WriteString("\n```\n")
	}
	if event.Response != "" {
		builder.WriteString("\n**Response**\n\n```http\n")
		// If the response is larger than 5 kb, truncate it before writing.
		if len(event.Response) > 5*1024 {
			builder.WriteString(event.Response[:5*1024])
			builder.WriteString(".... Truncated ....")
		} else {
			builder.WriteString(event.Response)
		}
		builder.WriteString("\n```\n\n")
	}

	if len(event.ExtractedResults) > 0 || len(event.Metadata) > 0 {
		builder.WriteString("**Extra Information**\n\n")
		if len(event.ExtractedResults) > 0 {
			builder.WriteString("**Extracted results**:\n\n")
			for _, v := range event.ExtractedResults {
				builder.WriteString("- ")
				builder.WriteString(v)
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}
		if len(event.Metadata) > 0 {
			builder.WriteString("**Metadata**:\n\n")
			for k, v := range event.Metadata {
				builder.WriteString("- ")
				builder.WriteString(k)
				builder.WriteString(": ")
				builder.WriteString(types.ToString(v))
				builder.WriteString("\n")
			}
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\n---\nGenerated by [Nuclei](https://github.com/projectdiscovery/nuclei)")
	data := builder.String()
	return data
}

// GetMatchedTemplate returns the matched template from a result event
func GetMatchedTemplate(event *output.ResultEvent) string {
	builder := &strings.Builder{}
	builder.WriteString(event.TemplateID)
	if event.MatcherName != "" {
		builder.WriteString(":")
		builder.WriteString(event.MatcherName)
	}
	if event.ExtractorName != "" {
		builder.WriteString(":")
		builder.WriteString(event.ExtractorName)
	}
	template := builder.String()
	return template
}
