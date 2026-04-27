package model

type AnalyticsSummary struct {
	Range             string  `json:"range"`
	Pageviews         int64   `json:"pageviews"`
	UniqueVisitors    int64   `json:"uniqueVisitors"`
	BounceRate        float64 `json:"bounceRate"`
	AvgSessionSeconds int64   `json:"avgSessionSeconds"`
	DeltaPageviews    float64 `json:"deltaPageviews"`
	PreviousPageviews int64   `json:"previousPageviews"`
}

type TopPage struct {
	Path    string `json:"path"`
	Views   int64  `json:"views"`
	Uniques int64  `json:"uniques"`
}

type TopReferrer struct {
	Referrer string `json:"referrer"`
	Views    int64  `json:"views"`
}

type TopCountry struct {
	Country string `json:"country"`
	Views   int64  `json:"views"`
}

type TimeseriesPoint struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type TopUtm struct {
	Value string `json:"value"`
	Views int64  `json:"views"`
}
