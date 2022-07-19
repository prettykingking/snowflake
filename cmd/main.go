package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/prettykingking/snowflake/pkg/cmd"
	"github.com/prettykingking/snowflake/pkg/config"
	"github.com/prettykingking/snowflake/pkg/provider"
	"github.com/prettykingking/snowflake/pkg/provider/aggregator"
	"github.com/prettykingking/snowflake/pkg/provider/grpc"
	"github.com/prettykingking/snowflake/pkg/safe"
	"github.com/prettykingking/snowflake/pkg/server"
)

func main() {
	cmd.RegisterCommand(cmd.Command{
		Name:  "run",
		Func:  cmdRun,
		Usage: "[--config <path>]",
		Short: "Starts the server process and blocks indefinitely",
		Long: `
Start the server process, optionally bootstrapped with an initial config file,
and blocks indefinitely until the server is stopped; i.e. runs server in
daemon mode (foreground).

If a config file is specified, it will be applied immediately after the process
is running.
`,
		Flags: func() *flag.FlagSet {
			fl := flag.NewFlagSet("run", flag.ExitOnError)
			fl.String("config", "", "Configuration file")
			return fl
		}(),
	})

	cmd.Execute()
}

// cmdRun runs services
func cmdRun(_ cmd.Flags, cf *config.Configuration) (int, error) {
	svr, err := setupServer(cf)
	if err != nil {
		return 1, err
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	svr.Start(ctx)
	defer svr.Close()

	svr.Wait()

	return 0, nil
}

// healthCheck check application environment
// include: network connections, configurations
func setupServer(cf *config.Configuration) (*server.Server, error) {
	pa := aggregator.NewProviderAggregator()

	ps := []provider.Provider{
		// more provider
		&grpc.Provider{},
	}

	for _, p := range ps {
		err := pa.AddProvider(p)
		if err != nil {
			return nil, err
		}
	}

	rp := safe.NewPool(context.Background())

	err := pa.Provide(cf, rp)
	if err != nil {
		return nil, err
	}

	return server.NewServer(rp), nil
}
