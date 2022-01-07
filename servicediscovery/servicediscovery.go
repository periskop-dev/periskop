package servicediscovery

import (
	"context"

	"log"

	gokit_log "github.com/go-kit/kit/log"
	"github.com/periskop-dev/periskop/config"
	prometheus_discovery "github.com/prometheus/prometheus/discovery"
	prometheus_discovery_config "github.com/prometheus/prometheus/discovery/config"
	prometheus_target_group "github.com/prometheus/prometheus/discovery/targetgroup"
	prometheus_labels "github.com/prometheus/prometheus/pkg/labels"
	prometheus_relabel "github.com/prometheus/prometheus/pkg/relabel"
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
	sdConfig       map[string]prometheus_discovery_config.ServiceDiscoveryConfig
	relabelConfigs []*prometheus_relabel.Config
}

func NewResolver(service config.Service) Resolver {
	sdConfig := map[string]prometheus_discovery_config.ServiceDiscoveryConfig{
		service.Name: service.ServiceDiscovery,
	}

	return Resolver{
		sdConfig:       sdConfig,
		relabelConfigs: service.RelabelConfigs,
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
			groups := <-manager.SyncCh()
			addresses := r.extractAddresses(groups)
			out <- ResolvedAddresses{
				Addresses: addresses,
			}
		}
	}()

	return out
}

func (r Resolver) extractAddresses(groups map[string][]*prometheus_target_group.Group) []string {
	var (
		addresses []string
		uniq      = make(map[string]struct{})
	)

	for _, groupArr := range groups {
		for i := 0; i < len(groupArr); i++ {
			group := groupArr[i]
			for _, target := range group.Targets {
				discoveredLabels := group.Labels.Merge(target)
				var labelMap = make(map[string]string)
				for k, v := range discoveredLabels.Clone() {
					labelMap[string(k)] = string(v)
				}

				processedLabels := prometheus_relabel.Process(prometheus_labels.FromMap(labelMap), r.relabelConfigs...)

				if processedLabels == nil {
					continue
				}

				// Deduplicate group targets with same address
				labels := processedLabels.Map()
				if _, prs := uniq[labels["__address__"]]; prs {
					continue
				}
				uniq[labels["__address__"]] = struct{}{}

				addresses = append(addresses, labels["__address__"])
			}
		}
	}
	return addresses
}
