package models

// Record represents a single row of data
type Record map[string]string

// Table represents the CSV data structure
type Table struct {
	Headers []string
	Records []Record
}

// QueryRequest represents the query parameters
type QueryRequest struct {
	Columns []string    `json:"columns"`
	Filter  FilterQuery `json:"filter"`
}

// FilterQuery represents filter parameters
type FilterQuery struct {
	Column   string `json:"column"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
