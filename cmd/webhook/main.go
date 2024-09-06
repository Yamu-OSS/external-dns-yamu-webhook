package main

import (
	"fmt"

	"github.com/Yamu-OSS/external-dns-yamu-webhook/cmd/webhook/init/configuration"
	"github.com/Yamu-OSS/external-dns-yamu-webhook/cmd/webhook/init/dnsprovider"
	"github.com/Yamu-OSS/external-dns-yamu-webhook/cmd/webhook/init/logging"
	"github.com/Yamu-OSS/external-dns-yamu-webhook/cmd/webhook/init/server"
	"github.com/Yamu-OSS/external-dns-yamu-webhook/pkg/webhook"
	log "github.com/sirupsen/logrus"
)

const banner = `
external-dns-provider-yamu
build tag: %s
build date: %s
git commit: %s
`

var (
	buildTime, gitCommitID, buildTag string
)

func main() {
	fmt.Printf(banner, buildTag, buildTime, gitCommitID)

	logging.Init()

	config := configuration.Init()
	provider, err := dnsprovider.Init(config)
	if err != nil {
		log.Fatalf("failed to initialize provider: %v", err)
	}

	main, health := server.Init(config, webhook.New(provider))
	server.ShutdownGracefully(main, health)
}
