package validate

import (
	"context"
	csvdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/csv"
)

// Validator is an interface that defines the methods that a service must implement.
type Validator interface {
	Schema(header, val string) error
	Filename(objects []collections.Object, csvMetadata csvdomain.CSVMetadata) error
	AssetID(objects []collections.Object, csvMetadata csvdomain.CSVMetadata, ctx context.Context) error
}
