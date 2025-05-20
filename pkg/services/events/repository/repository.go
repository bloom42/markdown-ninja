package repository

import (
	"time"

	"github.com/bloom42/stdx-go/memorycache"
)

type EventsRepository struct {
	cache *memorycache.Cache[string, any]
}

func NewEventsRepository() EventsRepository {
	cache := memorycache.New(
		memorycache.WithTTL[string, any](2 * time.Minute),
	)

	return EventsRepository{
		cache,
	}
}
