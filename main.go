package main

import (
	"flag"

	"github.com/comhttp/jorm/mod/cloudflare"
	"github.com/comhttp/jorm/mod/coin"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/utl"
	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
)

func main() {
	path := "/var/db/jorm/"
	c, _ := cfg.NewCFG(path, nil)
	config := &cfg.Config{}
	err := c.Read("conf", "conf", &config)
	utl.ErrorLog(err)

	loglevel := flag.String("loglevel", "debug", "Logging level (debug, info, warn, error)")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	switch *loglevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	cloudflare.CloudFlare(*config, &coin.CoinsShort{})

	log.Info().Msg("Starting CloudFlare API...")

	//log.Fatal(j.WWW.ListenAndServe())
}
