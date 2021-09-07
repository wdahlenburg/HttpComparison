package HttpComparison

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/wdahlenburg/raw_http_utils/utils"
)

type Similarity struct {
	Threshhold float64
}

func (s *Similarity) CompareRequests(baseline *http.Response, responses []*http.Response) ([]http.Response, error) {
	var differences []http.Response
	// Can ultimately perform this multithreaded with a wg

	rawBaseline := &utils.RawHttpResponse{
		Response: baseline,
	}

	result, err := rawBaseline.GetRawResponse()
	if err != nil {
		return nil, err
	}

	baseline_str := string(result)

	for _, response := range responses {
		weight := 1.0
		if baseline.StatusCode != response.StatusCode {
			weight = weight * 0.75 // Apply a 25% reduction if the status code does not match the baseline
		}

		raw := &utils.RawHttpResponse{
			Response: response,
		}

		res, err := raw.GetRawResponse()
		if err != nil {
			return nil, err
		}
		response_str := string(res)

		start := time.Now()
		hamming_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewHamming())
		duration := time.Since(start)
		fmt.Printf("Hamming: %.2f - %s\n", hamming_similarity, duration) // Output: 0.75

		// start = time.Now()
		// levenshtein_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewLevenshtein())
		// duration = time.Since(start)
		// fmt.Printf("Levenshtein: %.2f - %s\n", levenshtein_similarity, duration) // Output: 0.43

		start = time.Now()
		jaro_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewJaro())
		duration = time.Since(start)
		fmt.Printf("Jaro: %.2f - %s\n", jaro_similarity, duration) // Output: 0.78

		start = time.Now()
		jaro_winkler_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewJaroWinkler())
		duration = time.Since(start)
		fmt.Printf("Jaro-Winkler: %.2f - %s\n", jaro_winkler_similarity, duration) // Output: 0.80

		// start = time.Now()
		// swg_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewSmithWatermanGotoh())
		// duration = time.Since(start)
		// fmt.Printf("SWG: %.2f - %s\n", swg_similarity, duration) // Output: 0.82

		start = time.Now()
		sd_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewSorensenDice())
		duration = time.Since(start)
		fmt.Printf("Sorensen-Dice: %.2f - %s\n", sd_similarity, duration) // Output: 0.62

		start = time.Now()
		jaccard_similarity := strutil.Similarity(baseline_str, response_str, metrics.NewJaccard())
		duration = time.Since(start)
		fmt.Printf("Jaccard: %.2f - %s\n", jaccard_similarity, duration) // Output: 0.45

		sim := weight * jaccard_similarity
		if sim < s.Threshhold {
			fmt.Printf("Significant differentiation identified\n")
			differences = append(differences, *response)
		}

		start = time.Now()
		oc := metrics.NewOverlapCoefficient()
		similarity := strutil.Similarity(baseline_str, response_str, oc)
		duration = time.Since(start)
		fmt.Printf("Overlap: %.2f - %s\n", similarity, duration) // Output: 0.67

		fmt.Printf("\n")
	}

	return differences, nil
}