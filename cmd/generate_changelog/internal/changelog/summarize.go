package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const DefaultSummarizeModel = "claude-sonnet-4-20250514"
const MinContentLength = 256 // Minimum content length to consider for summarization

// getSummarizeModel returns the model to use for AI summarization
func getSummarizeModel() string {
	if model := os.Getenv("FABRIC_CHANGELOG_SUMMARIZE_MODEL"); model != "" {
		return model
	}
	return DefaultSummarizeModel
}

// SummarizeVersionContent takes raw version content and returns AI-enhanced summary
func SummarizeVersionContent(content string) (string, error) {
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("no content to summarize")
	}
	if len(content) < MinContentLength {
		// If content is too brief, return it as is
		return content, nil
	}

	model := getSummarizeModel()

	prompt := `Summarize the changes extracted from Git commit logs in a concise, professional way.
Pay particular attention to the following rules:
- Preserve the PR headers verbatim to your summary.
- I REPEAT: Do not change the PR headers in any way. They contain links to the PRs and Author Profiles.
- Use bullet points for lists and key changes (rendered using "-")
- Focus on the main changes and improvements.
- Avoid unnecessary details or preamble.
- Keep it under 800 characters.
- Be brief. List only the 5 most important changes along with the PR information which should be kept intact.
- If the content is too brief or you do not see any PR headers, return the content as is.`

	cmd := exec.Command("fabric", "-m", model, prompt)
	cmd.Stdin = strings.NewReader(content)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("fabric command failed: %w", err)
	}

	summary := strings.TrimSpace(string(output))
	if summary == "" {
		return "", fmt.Errorf("fabric returned empty summary")
	}

	return summary, nil
}
