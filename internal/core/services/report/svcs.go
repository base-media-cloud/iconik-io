package report

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/base-media-cloud/pd-iconik-io-rd/internal/core/domain/iconik/assets/collections/contents"
)

// Svc is a struct which contains all the information to generate a report.
type Svc struct {
	w *csv.Writer
}

// New is a function that returns a new instance of the report Svc struct.
func New(w *csv.Writer) *Svc {
	return &Svc{
		w: w,
	}
}

// Write is a method that writes data to a csv file and returns a byte stream.
func (svc *Svc) Write(ctx context.Context, objects []contents.ObjectDTO) error {
	var output []*contents.Object

	for j := range objects {
		if objects[j].ObjectType == "collections" {
			fmt.Printf("\nfound collection %s, collection id %s", objects[j].Title, objects[j].ID)
			if err := i.ProcessColl(objects[j].ID, 1, w); err != nil {
				return err
			}
			continue
		}
		output = append(output, objects[j])
	}

	toWrite, err := i.FormatObjects(output)
	if err != nil {
		return err
	}

	if err = svc.w.WriteAll(toWrite); err != nil {
		return err
	}

	return nil
}
