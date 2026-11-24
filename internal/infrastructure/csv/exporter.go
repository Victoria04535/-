package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"pos/internal/application"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Export(
	results []application.StudentResult,
	writer io.Writer,
) error {
	if writer == nil {
		return fmt.Errorf("export CSV: writer is nil")
	}
	csvWriter := csv.NewWriter(writer)
	header := []string{
		"student_id",
		"student_name",
		"total_score",
	}

	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("write CSV header: %w", err)
	}

	for _, result := range results {
		total := strconv.FormatFloat(
			result.Total,
			'f',
			2,
			64,
		)

		record := []string{
			result.Student.ID(),
			result.Student.Name(),
			total,
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("write CSV record: %w", err)
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("flush CSV: %w", err)
	}

	return nil
}
