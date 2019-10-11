package servicediscovery

import (
	"context"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/dns"
	"github.com/prometheus/prometheus/discovery/targetgroup"
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

type SRVResolver struct {
	dnsConfig dns.SDConfig
}

func NewResolver(c config.Service) SRVResolver {
	d, err := model.ParseDuration(c.ServiceDiscovery.RefreshInterval.String())
	if err != nil {
		panic(err)
	}
	names := []string{c.ServiceDiscovery.Name}
	dnsConfig := dns.SDConfig{
		Names:           names,
		RefreshInterval: d,
		Type:            c.ServiceDiscovery.Type,
		Port:            c.ServiceDiscovery.Port,
	}
	return SRVResolver{
		dnsConfig: dnsConfig,
	}
}

func (r SRVResolver) Resolve() <-chan ResolvedAddresses {
	out := make(chan ResolvedAddresses)

	srvDiscovery := dns.NewDiscovery(r.dnsConfig, nil)
	ctx := context.Background()
	groups := make(chan []*targetgroup.Group)
	go srvDiscovery.Run(ctx, groups)

	go func() {
		for {
			var addresses []string
			groupArr := <-groups
			for i := 0; i < len(groupArr); i++ {
				group := groupArr[i]
				for _, target := range group.Targets {
					addresses = append(addresses, string(target["__address__"]))
				}
			}
			out <- ResolvedAddresses{
				Addresses: addresses,
			}
		}
	}()

	return out
}
