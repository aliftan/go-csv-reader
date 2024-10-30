package services

import (
	"csv-reader/models"
	"encoding/csv"
	"io"
	"strings"
)

type CSVService struct {
	currentTable *models.Table
}

func NewCSVService() *CSVService {
	return &CSVService{
		currentTable: &models.Table{},
	}
}

func (s *CSVService) ProcessCSVFile(file io.Reader) error {
	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		return err
	}

	var records []models.Record
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		record := make(models.Record)
		for i, value := range row {
			record[headers[i]] = value
		}
		records = append(records, record)
	}

	s.currentTable = &models.Table{
		Headers: headers,
		Records: records,
	}

	return nil
}

func (s *CSVService) GetCurrentTable() *models.Table {
	return s.currentTable
}

func (s *CSVService) QueryRecords(columns []string, filter models.FilterQuery) []models.Record {
	filterFunc := func(record models.Record) bool {
		if filter.Value == "" {
			return true
		}

		value := record[filter.Column]
		switch filter.Operator {
		case "equals":
			return value == filter.Value
		case "contains":
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter.Value))
		case "greater":
			return value > filter.Value
		case "less":
			return value < filter.Value
		default:
			return true
		}
	}

	return s.SelectRecords(columns, filterFunc)
}

func (s *CSVService) SelectRecords(columns []string, where func(models.Record) bool) []models.Record {
	var result []models.Record

	for _, record := range s.currentTable.Records {
		if where == nil || where(record) {
			if len(columns) == 0 {
				result = append(result, record)
				continue
			}

			selectedRecord := make(models.Record)
			for _, col := range columns {
				if val, exists := record[col]; exists {
					selectedRecord[col] = val
				}
			}
			result = append(result, selectedRecord)
		}
	}

	return result
}
