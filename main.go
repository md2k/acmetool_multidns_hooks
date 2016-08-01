package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/md2k/acmetool_multidns_hooks/config"
	"github.com/md2k/acmetool_multidns_hooks/logging"
	"github.com/md2k/acmetool_multidns_hooks/providers/dns/cloudflare"
	"github.com/md2k/acmetool_multidns_hooks/providers/dns/route53"
	"github.com/md2k/acmetool_multidns_hooks/utils"
)

type Client struct {
	Acme     Acme
	Cfg      *config.DomainDetails
	Provider Provider
}

type Acme struct {
	Event     string
	Domain    string
	Ch_name   string
	Challenge string
}

type Provider interface {
	SetUp(domain, targetname, challenge string) error
	CleanUp(domain, targetname, challenge string) error
}

// var Hook Client

func main() {
	flag.Parse()
	args := flag.Args()

	// Init Color :D Logger
	logging.Init()

	// Get AcmeTool Path
	acme_path := os.Getenv("ACME_STATE_DIR")

	cfg, err := config.InitConfig(acme_path)
	if err != nil {
		logging.Log.Criticalf("%s", err.Error())
		os.Exit(1)
	}

	// We want process only DNS Events as we have functionality on as DNS Hook
	// switch event {
	// case "challenge-dns-start", "challenge-dns-stop":
	// 	logging.Log.Debugf("Event: %s | Domain: %s | Challenge File: %s | Challenge: %s", event, domain, challenge_filename, challenge)
	// 	// Get Domain Provider from config
	// 	// dns.PrepareChallenge(domain, challenge_filename, challenge)
	// default:
	// 	logging.Log.Criticalf("Unsupported event: %s", event)
	// 	os.Exit(1)
	// }

	switch args[0] {
	case "challenge-dns-start":
		logging.Log.Debugf("Event: %s | Domain: %s | Challenge File: %s | Challenge: %s", args[0], args[1], args[2], args[3])

		c, err := initDnsHook(cfg, args)
		if err != nil {
			logging.Log.Criticalf("%s", err.Error())
			os.Exit(1)
		}

		err = c.Provider.SetUp(c.Acme.Domain, c.Acme.Ch_name, c.Acme.Challenge)
		if err != nil {
			logging.Log.Criticalf("%s", err.Error())
			os.Exit(1)
		}
		logging.Log.Infof("Challenge record for domain %s successfully deployed to %s DNS Provider.", c.Acme.Domain, c.Cfg.Provider)
		os.Exit(0)

	case "challenge-dns-stop":
		logging.Log.Debugf("Event: %s | Domain: %s | Challenge File: %s | Challenge: %s", args[0], args[1], args[2], args[3])

		c, err := initDnsHook(cfg, args)
		if err != nil {
			logging.Log.Criticalf("%s", err.Error())
			os.Exit(1)
		}

		err = c.Provider.CleanUp(c.Acme.Domain, c.Acme.Ch_name, c.Acme.Challenge)
		if err != nil {
			logging.Log.Criticalf("%s", err.Error())
			os.Exit(1)
		}
		logging.Log.Infof("Challenge record for domain %s successfully removed from %s DNS Provider.", c.Acme.Domain, c.Cfg.Provider)
		os.Exit(0)
	case "live-updated":
		logging.Log.Debugf("Event: %s.", args[0])
		os.Exit(1)
	default:
		logging.Log.Criticalf("%s", err.Error())
		os.Exit(1)
	}
}

func initDnsHook(cfg *config.HookConfig, args []string) (*Client, error) {

	providerCfg, err := FindDomainProvider(cfg, args[1])
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	switch providerCfg.Provider {
	case "cloudflare":
		provider, err := cloudflare.InitProvider(providerCfg.AuthData["email"], providerCfg.AuthData["api_key"])
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		return &Client{
			Cfg: providerCfg,
			Acme: Acme{
				Event:     args[0],
				Domain:    args[1],
				Ch_name:   args[2],
				Challenge: args[3],
			},
			Provider: provider,
		}, nil
	case "route53":
		provider, err := route53.InitProvider(providerCfg.AuthData["key_id"], providerCfg.AuthData["access_key"])
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		fmt.Println(provider)
		return &Client{
			Cfg: providerCfg,
			Acme: Acme{
				Event:     args[0],
				Domain:    args[1],
				Ch_name:   args[2],
				Challenge: args[3],
			},
			Provider: provider,
		}, nil
	default:
		return nil, fmt.Errorf("%s", "No Provider's Module Available")
	}

}

func FindDomainProvider(cfg *config.HookConfig, domain string) (details *config.DomainDetails, err error) {
	// Find DNS Provider and Account Details for Specified Domain
	for provider, pval := range cfg.Providers {
		for account, aval := range pval.Accounts {
			if utils.StirngInSlice(domain, aval.Domains) {
				// fmt.Println(domain, provider, account, aval.AuthData)
				details = &config.DomainDetails{
					Provider: provider,
					Account:  account,
					AuthData: aval.AuthData,
				}
				return details, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Domain '%s' do not exist in configuration data", domain))
}
