package validate

import (
	"context"
	csvdomain "github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/csv"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collection"
)

// Validator is an interface that defines the methods that a service must implement.
type Validator interface {
	Schema(header, val string) error
	Filename(objects []collection.Object, csvMetadata csvdomain.CSVMetadata) error
	AssetID(objects []collection.Object, csvMetadata csvdomain.CSVMetadata, ctx context.Context) error
}
