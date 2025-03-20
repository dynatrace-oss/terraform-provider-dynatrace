package settings20

import (
	"log"
	"testing"
)

func TestSendClassicClientRequest(t *testing.T) {
	err := SendClassicClientRequest()

	if err != nil {
		log.Fatal(err)
	}
}

func TestSendPlatformClientRequest(t *testing.T) {
	err := SendPlatformClientRequest()

	if err != nil {
		log.Fatal(err)
	}
}
