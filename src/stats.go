package upgradebot

import (
	"math"
	"sort"
	"strings"
)

type PullRequestStats struct {
	Data PullRequestData

	FilesAddedCount    int
	FilesRemovedCount  int
	FilesModifiedCount int

	LinesAddedCount   int
	LinesRemovedCount int

	TopFilesChanged    []File
	TopPackagesChanged []PackageStats

	Assessment Assessment
}

type Assessment string

const (
	Good     Assessment = "Good"
	Warning             = "Warning"
	Conflict            = "Conflict"
)

type PackageStats struct {
	Name  string
	Count int
}

type ChangedFileStats struct {
	AssociatedPRs []PullRequestData
	File          File
	Assessment    Assessment
}

type Analysis struct {
	prStats   []PullRequestStats
	fileStats []ChangedFileStats
}

func GetAnalysis(tagCompare TagCompare, filesChangedByQuorum []string, expectedFileConflicts []string) Analysis {
	analysis := Analysis{}
	analysis.prStats = make([]PullRequestStats, len(tagCompare.PullRequests))

	// pre-processing
	mapFileAssessment := make(map[string]Assessment)
	for _, file := range filesChangedByQuorum {
		mapFileAssessment[file] = Warning
	}
	for _, file := range expectedFileConflicts {
		mapFileAssessment[file] = Conflict
	}

	// processing & ordering PRs
	for i, pr := range tagCompare.PullRequests {
		analysis.prStats[i] = getPullRequestStats(pr, mapFileAssessment)
	}

	sort.SliceStable(analysis.prStats, func(i, j int) bool {
		return analysis.prStats[i].Data.ClosedAt < analysis.prStats[j].Data.ClosedAt
	})

	analysis.fileStats = getChangedFilesStats(tagCompare, mapFileAssessment)

	return analysis
}

func getPullRequestStats(pr PullRequest, mapFileAssessment map[string]Assessment) PullRequestStats {
	stats := PullRequestStats{}

	stats.Data = pr.Data

	mapPackageChanged := make(map[string]int)

	for _, file := range pr.Files {
		stats.Assessment = Good
		if val, ok := mapFileAssessment[file.Filename]; stats.Assessment == Good && ok {
			if stats.Assessment != Conflict {
				stats.Assessment = val
			}
		}

		stats.LinesAddedCount += file.Additions
		stats.LinesRemovedCount += file.Deletions

		lastIndex := strings.LastIndex(file.Filename, "/")
		packagePath := file.Filename
		if lastIndex > 0 {
			packagePath = file.Filename[0:lastIndex]
		}
		mapPackageChanged[packagePath] = mapPackageChanged[packagePath] + 1

		switch file.Status {
		case "added":
			stats.FilesAddedCount += 1
		case "modified":
			stats.FilesModifiedCount += 1
		default:
			stats.FilesRemovedCount += 1
		}
	}

	sort.SliceStable(pr.Files, func(i, j int) bool {
		return pr.Files[i].getTotalModifications() > pr.Files[j].getTotalModifications()
	})

	stats.TopFilesChanged = pr.Files[0:int(math.Min(float64(len(pr.Files)), 5))]

	stats.TopPackagesChanged = make([]PackageStats, len(mapPackageChanged))

	i := 0
	for k, v := range mapPackageChanged {
		stats.TopPackagesChanged[i] = PackageStats{
			Name:  k,
			Count: v,
		}
		i++
	}

	sort.SliceStable(stats.TopPackagesChanged, func(i, j int) bool {
		return stats.TopPackagesChanged[i].Count > stats.TopPackagesChanged[j].Count
	})

	return stats
}

func getChangedFilesStats(tagCompare TagCompare, mapFileAssessment map[string]Assessment) []ChangedFileStats {
	prsPerFile := make(map[string][]PullRequestData)
	filePerFile := make(map[string]File)

	for _, file := range tagCompare.Files {
		prsPerFile[file.Filename] = make([]PullRequestData, 0)
		filePerFile[file.Filename] = file
	}

	for _, pr := range tagCompare.PullRequests {
		for _, file := range pr.Files {
			prsPerFile[file.Filename] = append(prsPerFile[file.Filename], pr.Data)
		}
	}

	stats := make([]ChangedFileStats, len(prsPerFile))

	i := 0
	for name, v := range prsPerFile {
		stats[i] = ChangedFileStats{AssociatedPRs: v, File: filePerFile[name], Assessment: mapFileAssessment[name]}
		i++
	}

	sort.SliceStable(stats, func(i, j int) bool {
		return stats[i].File.getTotalModifications() > stats[j].File.getTotalModifications()
	})

	return stats
}