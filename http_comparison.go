package HttpComparison

import (
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

type Similarity struct {
	Threshhold float64
}

func (s *Similarity) CompareStrResponses(baseline_str string, responses []string) ([]string, error) {
	var differences []string
	var baselineStatus string

	if len(baseline_str) != 0 {
		baselineStatus = strings.Split(strings.Split(baseline_str, "\n")[0], " ")[1]
	} else {
		baselineStatus = ""
	}

	for _, response := range responses {
		weight := 1.0
		responseCode := strings.Split(strings.Split(response, "\n")[0], " ")[1]
		if baselineStatus != responseCode {
			weight = weight * 0.75 // Apply a 25% reduction if the status code does not match the baseline
		}

		jaccard_similarity := strutil.Similarity(baseline_str, response, metrics.NewJaccard())

		sim := weight * jaccard_similarity
		if sim < s.Threshhold {
			differences = append(differences, response)
		}
	}

	return differences, nil
}
