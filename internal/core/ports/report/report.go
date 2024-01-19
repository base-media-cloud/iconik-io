package report

//go:generate mockgen -source report.go -destination=../../../mocks/report_mocks/report.go -package=report_mocks

import (
	"context"
)

// Reporter is an interface that defines the methods a report can implement.
type Reporter interface {
	Write(ctx context.Context, data interface{}) ([]byte, error)
}
