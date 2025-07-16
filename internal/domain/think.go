package domain

import "regexp"

// StripThinkBlocks removes any content between the provided start and end tags
// from the input string. Whitespace following the end tag is also removed so
// output resumes at the next non-empty line.
func StripThinkBlocks(input, startTag, endTag string) string {
	if startTag == "" || endTag == "" {
		return input
	}
	pattern := "(?s)" + regexp.QuoteMeta(startTag) + ".*?" + regexp.QuoteMeta(endTag) + "\\s*"
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, "")
}
