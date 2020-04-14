package servicediscovery

import (
	"context"

	"log"

	gokit_log "github.com/go-kit/kit/log"
	prometheus_discovery "github.com/modularise/prometheus-discovery/discovery"
	prometheus_discovery_config "github.com/modularise/prometheus-discovery/discovery/config"
	"github.com/soundcloud/periskop/config"
)

type ResolvedAddresses struct {
	Addresses []string
}

func EmptyResolvedAddresses() ResolvedAddresses {
	return ResolvedAddresses{
		Addresses: make([]string, 0),
	}
}

type Resolver struct {
	sdConfig map[string]prometheus_discovery_config.ServiceDiscoveryConfig
}

func NewResolver(service config.Service) Resolver {
	sdConfig := map[string]prometheus_discovery_config.ServiceDiscoveryConfig{
		service.Name: service.ServiceDiscovery,
	}

	return Resolver{
		sdConfig: sdConfig,
	}
}

func (r Resolver) Resolve() <-chan ResolvedAddresses {
	ctx := context.Background()
	out := make(chan ResolvedAddresses)
	manager := prometheus_discovery.NewManager(ctx, gokit_log.NewNopLogger())

	err := manager.ApplyConfig(r.sdConfig)
	if err != nil {
		log.Fatal("Could not apply SD configuration")
	}

	go func() {
		err = manager.Run()
	}()

	if err != nil {
		log.Fatal("Could not initialize SD manager")
	}

	go func() {
		for {
			var addresses []string
			groups := <-manager.SyncCh()
			for _, groupArr := range groups {
				for i := 0; i < len(groupArr); i++ {
					group := groupArr[i]
					for _, target := range group.Targets {
						addresses = append(addresses, string(target["__address__"]))
					}
				}
			}
			out <- ResolvedAddresses{
				Addresses: addresses,
			}
		}
	}()

	return out
}
