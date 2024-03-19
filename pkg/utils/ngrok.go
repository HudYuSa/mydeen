package utils

import (
	"context"
	"log"

	"golang.ngrok.com/ngrok"
	ngrokConfig "golang.ngrok.com/ngrok/config"
)

func RunNgrok() (ngrok.Tunnel, error) {
	// ngrok
	tun, err := ngrok.Listen(context.Background(),
		ngrokConfig.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	log.Println("tunnel created:", tun.URL())

	return tun, err
}
