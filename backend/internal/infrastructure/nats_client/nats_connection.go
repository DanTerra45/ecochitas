package nats_client

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

func New_nats_connection(nats_url string) (*nats.Conn, error) {
	nats_connection, connect_nats_error := nats.Connect(
		nats_url,
		nats.Name("ecochitas_backend"),
		nats.Timeout(10*time.Second),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
	)
	if connect_nats_error != nil {
		return nil, fmt.Errorf("failed_to_connect_nats: %w", connect_nats_error)
	}

	return nats_connection, nil
}
