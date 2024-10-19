package main

import "context"

// fetcher works concurrently with main go route.
// it will continuously fetch the API and update the data on the given interval
func (c *Config) dataSourceFetcher(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	}
}
