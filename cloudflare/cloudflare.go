package cloudflare

import (
	"context"
	"github.com/comhttp/jorm/mod/coin"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/utl"
	"github.com/rs/zerolog/log"

	cf "github.com/cloudflare/cloudflare-go"
)

func CloudFlare(c cfg.Config, cfCoins *coin.CoinsShort) {
	//log.Print("CONFIGCONFIGCONFIGCONFIGCONFIGCONFIGCONFIG", cfg.C)
	ctx := context.Background()
	// Construct a new API object
	api, err := cf.NewWithAPIToken(c.CF.CloudFlareAPItoken)
	utl.ErrorLog(err)
	for _, tld := range c.COMHTTP {
		createDNS(cfCoins, api, ctx, "com-http."+tld)
		//delAllCNameDNS(api, ctx, "com-http."+tld)
	}
	//createDNS(j,api, ctx, "com-http.us")
	//delAllCNameDNS(api, ctx, "com-http.us")
}

func createDNS(cfCoins *coin.CoinsShort, api *cf.API, ctx context.Context, domain string) {
	// Fetch the zone ID
	id, err := api.ZoneIDByName(domain) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal().Err(err)
	}
	// Fetch all records for a zone
	recs, err := api.DNSRecords(context.Background(), id, cf.DNSRecord{})
	if err != nil {
		log.Fatal().Err(err)
	}
	var registrated []string
	for _, r := range recs {
		if r.Type == "CNAME" {
			registrated = append(registrated, r.Name)
		}
	}
	for _, cfCoin := range cfCoins.C {
		//_, err := http.Get("https://" + slug + "." + domain)
		//if err != nil {
		setDNS(api, ctx, registrated, domain, cfCoin.Slug)
	}
}

func setDNS(api *cf.API, ctx context.Context, registrated []string, domain, slug string) {
	var exist bool
	for _, reg := range registrated {
		if slug+"."+domain == reg {
			log.Print("Ima:", slug+"."+domain)
			exist = true
		} else {
			exist = false
		}
	}
	if !exist {
		id, err := api.ZoneIDByName(domain)
		utl.ErrorLog(err)
		t := true
		_, err = api.CreateDNSRecord(ctx, id, cf.DNSRecord{
			Type:    "CNAME",
			Name:    slug,
			Content: domain,
			TTL:     1,
			Proxied: &t,
		})
		utl.ErrorLog(err)
		log.Print("Created subdomain: ", slug+"."+domain)
	}
}

func delAllCNameDNS(api *cf.API, ctx context.Context, domain string) {
	// Fetch the zone ID
	id, err := api.ZoneIDByName(domain) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal().Err(err)
	}
	// Fetch all records for a zone
	recs, err := api.DNSRecords(context.Background(), id, cf.DNSRecord{})
	if err != nil {
		log.Fatal().Err(err)
	}
	for _, r := range recs {
		if r.Type == "CNAME" {
			go delDNS(api, ctx, id, r.ID)
		}
	}
}

func delDNS(api *cf.API, ctx context.Context, zoneId, id string) {
	err := api.DeleteDNSRecord(ctx, zoneId, id)
	utl.ErrorLog(err)
	log.Print("DeleteDNSRecord rrrrr:", id)
}
