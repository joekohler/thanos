// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package storecache

import (
	"context"
	"fmt"

	"github.com/oklog/ulid"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"golang.org/x/exp/slices"
)

type FilteredIndexCache struct {
	cache        IndexCache
	enabledItems []string
}

// NewFilteredIndexCache creates a filtered index cache based on enabled items.
func NewFilteredIndexCache(cache IndexCache, enabledItems []string) *FilteredIndexCache {
	return &FilteredIndexCache{
		cache:        cache,
		enabledItems: enabledItems,
	}
}

// StorePostings sets the postings identified by the ulid and label to the value v,
// if the postings already exists in the cache it is not mutated.
func (c *FilteredIndexCache) StorePostings(blockID ulid.ULID, l labels.Label, v []byte, tenant string) {
	if len(c.enabledItems) == 0 || slices.Contains(c.enabledItems, cacheTypePostings) {
		c.cache.StorePostings(blockID, l, v, tenant)
	}
}

// FetchMultiPostings fetches multiple postings - each identified by a label -
// and returns a map containing cache hits, along with a list of missing keys.
func (c *FilteredIndexCache) FetchMultiPostings(ctx context.Context, blockID ulid.ULID, keys []labels.Label, tenant string) (hits map[labels.Label][]byte, misses []labels.Label) {
	if len(c.enabledItems) == 0 || slices.Contains(c.enabledItems, cacheTypePostings) {
		return c.cache.FetchMultiPostings(ctx, blockID, keys, tenant)
	}
	return nil, keys
}

// StoreExpandedPostings stores expanded postings for a set of label matchers.
func (c *FilteredIndexCache) StoreExpandedPostings(blockID ulid.ULID, matchers []*labels.Matcher, v []byte, tenant string) {
	if len(c.enabledItems) == 0 || slices.Contains(c.enabledItems, cacheTypeExpandedPostings) {
		c.cache.StoreExpandedPostings(blockID, matchers, v, tenant)
	}
}

// FetchExpandedPostings fetches expanded postings and returns cached data and a boolean value representing whether it is a cache hit or not.
func (c *FilteredIndexCache) FetchExpandedPostings(ctx context.Context, blockID ulid.ULID, matchers []*labels.Matcher, tenant string) ([]byte, bool) {
	if len(c.enabledItems) == 0 || slices.Contains(c.enabledItems, cacheTypeExpandedPostings) {
		return c.cache.FetchExpandedPostings(ctx, blockID, matchers, tenant)
	}
	return nil, false
}

// StoreSeries sets the series identified by the ulid and id to the value v,
// if the series already exists in the cache it is not mutated.
func (c *FilteredIndexCache) StoreSeries(blockID ulid.ULID, id storage.SeriesRef, v []byte, tenant string) {
	if len(c.enabledItems) == 0 || slices.Contains(c.enabledItems, cacheTypeSeries) {
		c.cache.StoreSeries(blockID, id, v, tenant)
	}
}

// FetchMultiSeries fetches multiple series - each identified by ID - from the cache
// and returns a map containing cache hits, along with a list of missing IDs.
func (c *FilteredIndexCache) FetchMultiSeries(ctx context.Context, blockID ulid.ULID, ids []storage.SeriesRef, tenant string) (hits map[storage.SeriesRef][]byte, misses []storage.SeriesRef) {
	if len(c.enabledItems) == 0 || slices.Contains(c.enabledItems, cacheTypeSeries) {
		return c.cache.FetchMultiSeries(ctx, blockID, ids, tenant)
	}
	return nil, ids
}

func ValidateEnabledItems(enabledItems []string) error {
	for _, item := range enabledItems {
		switch item {
		// valid
		case cacheTypePostings, cacheTypeExpandedPostings, cacheTypeSeries:
		default:
			return fmt.Errorf("unsupported item type %s", item)
		}
	}
	return nil
}
