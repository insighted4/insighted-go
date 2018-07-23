package pgstore

import (
	"context"
	"time"

	"github.com/go-pg/pg"
	"github.com/insighted4/insighted-go/eventsourcing"
	"github.com/sirupsen/logrus"
)

// Postgres provides a Postgres backed store
type Postgres struct {
	db     *pg.DB
	logger logrus.FieldLogger
}

// Save implements the Store interface and saves records, serialized events, in Postgres
func (p *Postgres) Save(ctx context.Context, aggregateID string, records ...eventsourcing.Record) error {
	return nil
}

// Load implements the Store interface and retrieve events from Postgres
func (p *Postgres) Load(ctx context.Context, aggregateID string, fromVersion, toVersion int) (eventsourcing.History, error) {
	return nil, nil
}

// New returns a Postgres backed store
func New(options *pg.Options, logger logrus.FieldLogger) *Postgres {
	logger = logger.WithField("component", "Postgres")

	db := pg.Connect(options)
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		logger.WithField("query", query).
			WithField("latency", time.Since(event.StartTime)).
			Debugf("Query processed")
	})

	return &Postgres{
		db:     db,
		logger: logger,
	}
}
