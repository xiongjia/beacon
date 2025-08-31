package storage

import (
	"github.com/xiongjia/beacon/pkg/metric"
)

type (
	storage struct {
	}

	StorageOption func(*storage)

	Storage interface {
		Append(rcd []metric.Record) error

		Select(metric string, labels []metric.Label, start, end int64) (points []*metric.DataPoint, err error)

		Close() error
	}
)

func NewStorage(opts ...StorageOption) (Storage, error) {
	// xxx
	return nil, nil
}

func (s *storage) Append(rcd []metric.Record) error {
	return nil
}
