// Command wikipedia_api is a small client for the Wikimedia Pageviews REST API.
//
// It demonstrates four things against the public "metrics/pageviews" endpoints:
//
//  1. getTopArticles          — the top-N most viewed articles for a single day.
//  2. calcPercentage          — each country's share of total pageviews for a month.
//  3. calcEstimatedViews      — a rough split of a day's top-article views across countries.
//  4. topArticlesOverThreshold — a concurrent scan of every day in a month to find
//     articles that stayed in the top-N for more than a given number of days.
//
// The API requires a descriptive User-Agent header on every request (see userAgent
// below); requests without one are rejected. No API key is needed — the endpoints
// are public and read-only.
//
// API reference: https://wikimedia.org/api/rest_v1/#/Pageviews%20data
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	// project is the Wikimedia project to query (English Wikipedia here).
	project = "en.wikipedia"
	// access filters by access method: "all-access", "desktop", "mobile-app", or "mobile-web".
	access = "all-access"
	// baseURL is the root of the Wikimedia Pageviews REST API.
	baseURL = "https://wikimedia.org/api/rest_v1/metrics/pageviews"
	// userAgent identifies this client to Wikimedia; a descriptive UA is required by
	// their policy, and requests without one are throttled or rejected.
	userAgent = "app/version (abc@example.com)"
)

// Data is the top-level response from the "top articles per day" endpoint.
type Data struct {
	Items []Item `json:"items"`
}

// Item wraps the list of ranked articles inside a Data response.
type Item struct {
	Articles []Article `json:"articles"`
}

// Article is a single ranked article with its view count for the requested day.
type Article struct {
	Article string `json:"article"` // article title (underscores for spaces, e.g. "Main_Page")
	Views   int    `json:"views"`   // total pageviews for the day
	Rank    int    `json:"rank"`    // 1-based position in the top list
}

// CountryData is the top-level response from the "top viewers by country" endpoint.
type CountryData struct {
	Items []CountryItem `json:"items"`
}

// CountryItem wraps the list of per-country view stats inside a CountryData response.
type CountryItem struct {
	Countries []ViewsByCountry `json:"countries"`
}

// ViewsByCountry holds a single country's pageview stats for a month.
//
// Note: the API returns Views as a bucketed string (e.g. "1000-9999") for privacy,
// so ViewsCeil — the numeric upper bound of that bucket — is the field used for any
// arithmetic in this program.
type ViewsByCountry struct {
	Country   string `json:"country"`    // ISO country code, e.g. "US"
	Views     string `json:"views"`      // privacy-bucketed range as a string
	Rank      int    `json:"rank"`       // 1-based position among countries
	ViewsCeil int    `json:"views_ceil"` // numeric ceiling of the Views bucket
}

// client is a shared HTTP client with a timeout so no single request hangs forever.
var client = &http.Client{Timeout: 20 * time.Second}

func main() {

	month := 04
	year := 2022
	day := 01

	// Get top 5 articles
	getTopArticles(5, year, month, day)
	// Calculate the top 10 coutry view %
	calcPercentage(10, year, month)

	calcEstimatedViews(2022, 04, 01)

	topArticlesOverThresholdConcurrent(2022, 4, 30, 10, 3)
}

// getArticles fetches the ranked list of most-viewed articles for a single day.
// It returns the articles from the first (and only) item in the response, or an
// error if the request fails, the status is non-200, or the day has no data.
func getArticles(year, month, day int) ([]Article, error) {
	url := fmt.Sprintf("%s/top/%s/%s/%d/%02d/%02d", baseURL, project, access, year, month, day)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET error for %s returned %d", url, resp.StatusCode)
	}

	var r Data
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if len(r.Items) == 0 {
		return nil, fmt.Errorf("No records for year %d, month %d, day %d", year, month, day)
	}

	return r.Items[0].Articles, nil
}

// getViewsByCountry fetches the per-country pageview breakdown for a whole month.
// It returns the country list from the first item in the response, or an error if
// the request fails, the status is non-200, or the month has no data.
func getViewsByCountry(year, month int) ([]ViewsByCountry, error) {
	url := fmt.Sprintf("%s/top-by-country/%s/%s/%d/%02d", baseURL, project, access, year, month)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET error for %s returned %d", url, resp.StatusCode)
	}

	// // ---- ADD THESE TWO LINES ----
	// raw, _ := io.ReadAll(resp.Body)
	// fmt.Printf("RAW: %s\n", raw)
	// // ------------------------------

	var r CountryData
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if len(r.Items) == 0 {
		return nil, fmt.Errorf("No country records for year %d, month %d", year, month)
	}

	return r.Items[0].Countries, nil
}

// getTopArticles prints the top-n articles for the given day. It truncates the list
// to n entries and exits the process (log.Fatalf) if the fetch fails.
func getTopArticles(n, year, month, day int) {
	articles, err := getArticles(year, month, day)
	if err != nil {
		log.Fatalf("error in getting articles for %02d-%02d-%d", day, month, year)
	}

	if len(articles) > n {
		articles = articles[:n]
	}

	fmt.Printf("Top %d articles for %02d-%02d-%d:\n", n, day, month, year)
	for _, a := range articles {
		fmt.Printf("Rank - %d Article - %s Views - %d\n", a.Rank, a.Article, a.Views)
	}
}

// calcPercentage prints, for the top-n countries, each country's share of the
// month's total pageviews. Percentages are computed against the summed ViewsCeil
// across all countries (not just the top n), so they represent a true fraction of
// the whole. It exits (log.Fatalf) on fetch error or if no views were recorded.
func calcPercentage(n, year, month int) {
	totalViews := 0
	data, err := getViewsByCountry(year, month)
	if err != nil {
		log.Fatalf("error in getting views by country for %02d-%d", month, year)
	}

	// Sum every country's bucket ceiling to get the denominator for the percentages.
	for _, viewData := range data {
		totalViews += viewData.ViewsCeil

	}

	if totalViews == 0 {
		log.Fatalf("No views recorded per country")
	}

	if n > len(data) {
		n = len(data)
	}

	for _, c := range data[:n] {
		pct := 100 * float64(c.ViewsCeil) / float64(totalViews)
		fmt.Printf("Country - %s, Rank - %d, Percentage views - %.2f \n", c.Country, c.Rank, pct)
	}

}

// calcEstimatedViews prints a rough estimate of how a day's top-article views split
// across countries. It multiplies the total views of the day's top articles by each
// country's share of monthly traffic — a back-of-envelope approximation, since the
// per-country data is monthly while the article data is for a single day.
func calcEstimatedViews(year int, month int, day int) {
	articleCount := 5
	countryCount := 10
	articleViews, err := getArticles(year, month, day)
	if err != nil {
		log.Fatalf("error in getting articles for %02d-%02d-%d", day, month, year)
	}

	countryViews, err := getViewsByCountry(year, month)
	if err != nil {
		log.Fatalf("error in getting views by country for %02d-%d", month, year)
	}

	// Denominator: total monthly views across every country (before truncation).
	totalCountryViews := 0
	for _, countryView := range countryViews {
		totalCountryViews += countryView.ViewsCeil
	}

	if len(articleViews) > articleCount {
		articleViews = articleViews[:articleCount]
	}
	if len(countryViews) > countryCount {
		countryViews = countryViews[:countryCount]
	}

	// Numerator basis: combined views of the day's top articles.
	articleViewsTotal := 0

	for _, article := range articleViews {
		articleViewsTotal += article.Views
	}

	// Distribute the article views by each country's share of monthly traffic.
	for _, country := range countryViews {
		estimatedViews := float64(country.ViewsCeil) / float64(totalCountryViews) * float64(articleViewsTotal)
		fmt.Printf("Country - %s, Rank - %d, Total Views - %d, Estimated top views - %.2f\n", country.Country, country.Rank, country.ViewsCeil, estimatedViews)
	}

}

// topArticlesOverThresholdConcurrent scans every day of a month in parallel and
// reports which articles appeared in the daily top-topN on more than minDays days.
//
// One goroutine is launched per day; each fetches that day's top articles and
// increments a shared counter per article title. A mutex guards the map because
// many goroutines write to it, and a WaitGroup blocks until all fetches complete.
// Days that fail to fetch are logged and skipped rather than aborting the whole run.
// Results are sorted by day-count descending, then by title for stable ordering.
func topArticlesOverThresholdConcurrent(year, month, daysInMonth, topN, minDays int) {
	dayCount := map[string]int{}
	var mu sync.Mutex     // guards dayCount — many goroutines write to it
	var wg sync.WaitGroup // waits for all day-fetches to finish

	for day := 1; day <= daysInMonth; day++ {
		wg.Add(1)
		go func(day int) { // pass day as arg to avoid loop-var capture
			defer wg.Done()

			articles, err := getArticles(year, month, day)
			if err != nil {
				log.Printf("skipping %02d-%02d-%d: %v", day, month, year, err)
				return
			}
			if len(articles) > topN {
				articles = articles[:topN]
			}

			mu.Lock()
			for _, a := range articles {
				dayCount[a.Article]++
			}
			mu.Unlock()
		}(day)
	}
	wg.Wait() // block until every goroutine finishes

	// Collect articles that cleared the minDays threshold into a sortable slice.
	type result struct {
		name string
		days int
	}
	var results []result
	for name, days := range dayCount {
		if days > minDays {
			results = append(results, result{name, days})
		}
	}
	// Sort by most days in the top list first; break ties alphabetically for determinism.
	sort.Slice(results, func(i, j int) bool {
		if results[i].days != results[j].days {
			return results[i].days > results[j].days
		}
		return results[i].name < results[j].name
	})

	fmt.Printf("Articles in the top %d for more than %d days in %02d-%d:\n", topN, minDays, month, year)
	for _, r := range results {
		fmt.Printf("  %-30s %d days\n", r.name, r.days)
	}
}
