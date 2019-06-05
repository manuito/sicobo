package store

type BookCollection struct {
	Name         string
	TotalVolumes int64
	Active       bool
}

type Book struct {
	Title            string
	Picture          string
	Isbn             string
	PublishedDate    string
	Category         string
	Snippet          string
	PageCount        int
	Authors          []Author
	Collection       BookCollection
	CandidateDetails CandidateDetails
}

type CandidateDetails struct {
	DegradedSource bool
	Collections    []string
	Titles         []string
	Pictures       []string
}

type Author struct {
	Name string
}
