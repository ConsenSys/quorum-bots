package upgradebot

import (
	"fmt"
	"strings"
)

func CreateMarkdownHeader() string {
	builder := strings.Builder{}

	builder.WriteString("## Actions\n\n")

	builder.WriteString("### Tasks to be done\n\n")

	builder.WriteString("- [ ] Assign someone responsible for this go-ethereum upstream upgrade to the PR\n")
	builder.WriteString("- [ ] Review new features and fixes on the Release Notes\n")
	builder.WriteString("- [ ] Review Pull Requests in the analysis️\n")
	builder.WriteString("- [ ] Review new features need to be updated or enhanced for Quorum\n")
	builder.WriteString("- [ ] Solve all conflicts\n")
	builder.WriteString("- [ ] Review if new unit tests or Acceptance Tests are required\n")
	builder.WriteString("- [ ] Document in the Extra Changes section any changes or new code that was added as part of this PR\n")

	builder.WriteString("\n\n")

	builder.WriteString("### Extra Changes\n\n")

	builder.WriteString("* **\\<Example\\>**: \\<change\\>\n")

	builder.WriteString("\n\n")

	return builder.String()
}

func CreateMarkdownReleaseSection(data ReleaseData) string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "## Summary of: %s\n\n", data.Name)

	fmt.Fprintf(&builder, "* Version: %s\n", data.Tag)
	fmt.Fprintf(&builder, "* Published: %s\n", data.PublishedAt)

	builder.WriteString("\n\n")

	builder.WriteString("### Release notes \n\n")

	builder.WriteString(data.Body)

	return builder.String()
}

func CreateMarkdownAnalysisSection(analysis Analysis) string {
	builder := strings.Builder{}

	builder.WriteString("## Analysis\n\n")

	fmt.Fprintf(&builder, "### Legend\n\n")

	builder.WriteString("File Stats: (A) Added, (M) Modified and (R) Removed\n\n")
	builder.WriteString("Line Stats: (A) Added and (R) Removed\n\n")

	builder.WriteString("Assessment:\n\n")
	builder.WriteString("* ✅ No conflict expected\n")
	builder.WriteString("* ⚠ Review required to assess changes\n")
	builder.WriteString("* ‼️ Conflicts expected and review required\n")

	builder.WriteString("\n\n")

	fmt.Fprintf(&builder, "### Summary of %d merged Pull Requests\n\n", len(analysis.prStats))

	builder.WriteString("\n\n")

	builder.WriteString("| 🔍 | Link | Title | File Stats<br>M/A/R | Packages changed<br>(files changed) | Line Stats<br>A/R | Top 5 Changed Files<br>(lines changed) |\n")
	builder.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- | :--- |\n")

	for _, stats := range analysis.prStats {
		fmt.Fprintf(&builder, "| %s | [#%d](%s) | ``%s`` | %s | %s | %s | %s |\n",
			getAssessmentEmoji(stats.Assessment),
			stats.Data.Number,
			stats.Data.HtmlUrl,
			stats.Data.Title,
			createMarkdownPullRequestFileStats(stats),
			createMarkdownPullRequestPackageChangedStats(stats),
			createMarkdownPullRequestLineStats(stats),
			createMarkdownPullRequestTopChangedStats(stats))
	}

	builder.WriteString("\n\n")

	builder.WriteString("### Summary of Changed files\n\n")

	builder.WriteString("| 🔍 | File | Lines Changed | Linked PR |\n")
	builder.WriteString("| :--- | :--- | :--- | :--- |\n")

	for _, stat := range analysis.fileStats {
		fmt.Fprintf(&builder, "| %s | ``%s`` | %d | %s |\n",
			getAssessmentEmoji(stat.Assessment),
			stat.File.Filename,
			stat.File.getTotalModifications(),
			createMarkdownPullRequestDataListStats(stat.AssociatedPRs))
	}

	builder.WriteString("\n\n")

	return builder.String()
}

func createMarkdownPullRequestDataListStats(prDataArray []PullRequestData) string {
	builder := strings.Builder{}

	for _, data := range prDataArray {
		fmt.Fprintf(&builder, "[#%d](%s)<br>", data.Number, data.HtmlUrl)
	}

	return builder.String()
}

func createMarkdownPullRequestFileStats(stats PullRequestStats) string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "%d/%d/%d<br>", stats.FilesModifiedCount, stats.FilesAddedCount, stats.FilesRemovedCount)

	return builder.String()
}

func createMarkdownPullRequestLineStats(stats PullRequestStats) string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "<span class=\"text-green\">%d</span>/<span class=\"text-red\">%d</span><br>", stats.LinesAddedCount, stats.LinesRemovedCount)

	return builder.String()
}

func createMarkdownPullRequestPackageChangedStats(stats PullRequestStats) string {
	builder := strings.Builder{}

	for _, f := range stats.TopPackagesChanged {
		if f.Count > 0 {
			fmt.Fprintf(&builder, "``%s`` (%d)<br>", f.Name, f.Count)
		}
	}

	return builder.String()
}

func getAssessmentEmoji(assessment Assessment) string {
	switch assessment {
	case Conflict:
		return "‼️"
	case Warning:
		return "⚠️"
	default:
		return "✅"
	}
}

func createMarkdownPullRequestTopChangedStats(stats PullRequestStats) string {
	builder := strings.Builder{}

	for _, f := range stats.TopFilesChanged {
		if f.Changes > 0 {
			fmt.Fprintf(&builder, "``%s`` (%d)<br>", f.Filename, f.getTotalModifications())
		}
	}

	return builder.String()
}