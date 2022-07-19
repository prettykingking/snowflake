package aggregator

import (
	"github.com/prettykingking/snowflake/pkg/config"
	"github.com/prettykingking/snowflake/pkg/provider"
	"github.com/prettykingking/snowflake/pkg/safe"
)

// ProviderAggregator aggregates providers.
type ProviderAggregator struct {
	providers []provider.Provider
}

// NewProviderAggregator returns an aggregate of empty providers.
func NewProviderAggregator() *ProviderAggregator {
	return &ProviderAggregator{}
}

// AddProvider adds a provider in the providers map.
func (pa *ProviderAggregator) AddProvider(p provider.Provider) error {
	err := p.Init()
	if err != nil {
		return err
	}

	pa.providers = append(pa.providers, p)

	return nil
}

// Provide calls the provide method of every providers.
func (pa *ProviderAggregator) Provide(cf *config.Configuration, rp *safe.Pool) error {
	for _, p := range pa.providers {
		err := p.Provide(cf, rp)
		if err != nil {
			return err
		}
	}

	return nil
}
