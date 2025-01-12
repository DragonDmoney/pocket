package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO(team): Define the configs we need and add more tests here.

func TestLoadConfigFromJson(t *testing.T) {
	config := `{
		"root_dir": "/go/src/github.com/pocket-network",
		"genesis_source": {
			"file": {
				"path": "build/config/genesis.json"
			}
		},
		"private_key": "2e00000000000000000000000000000000000000000000000000000000000000264a0707979e0d6691f74b055429b5f318d39c2883bb509310b67424252e9ef2",
		"pre2p": {
		  "consensus_port": 8080,
		  "use_raintree": true,
		  "connection_type": "tcp"
		},
		"p2p": {
		  "protocol": "tcp",
		  "address": "0.0.0.0:8081",
		  "external_ip": "172.18.0.1:8081",
		  "peers": [
			"172.18.0.1:8081",
			"172.18.0.1:8082",
			"172.18.0.1:8083",
			"172.18.0.1:8084"
		  ]
		},
		"consensus": {
		  "max_mempool_bytes": 500000000,
		  "max_block_bytes": 4000000,
		  "pacemaker": {
			"timeout_msec": 5000,
			"manual": true,
			"debug_time_between_steps_msec": 1000
		  }
		},
		"pre_persistence": {
		  "capacity": 99999,
		  "mempool_max_bytes": 99999,
		  "mempool_max_txs": 99999
		},
		"persistence": {
		  "postgres_url": "postgres://postgres:postgres@pocket-db:5432/postgres",
		  "schema": "node1"
		},
		"utility": {}
	  }`

	c := Config{}
	err := json.Unmarshal([]byte(config), &c)
	require.NoError(t, err)
}
