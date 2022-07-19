package provider

import (
	"github.com/prettykingking/snowflake/pkg/config"
	"github.com/prettykingking/snowflake/pkg/safe"
)

// Provider defines methods of a provider.
type Provider interface {
	// Init services before run
	Init() error

	// Provide run services in goroutines pool.
	Provide(cf *config.Configuration, rp *safe.Pool) error
}
