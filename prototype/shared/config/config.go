package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const DefaultGenesis = "build/config/genesis.json"

// TODO: Should we make an interface so we can mock it and inject it for testing?
type Config struct {
	RootDir    string `json:"root_dir"`
	PrivateKey string `json:"private_key"`
	Genesis    string `json:"genesis"`

	PREP2P      *PREP2PConfig      `json:"prep2p"`
	P2P         *P2PConfig         `json:"p2p"`
	Consensus   *ConsensusConfig   `json:"consensus"`
	Persistence *PersistenceConfig `json:"persistence"`
	Utility     *UtilityConfig     `json:"utility"`
}

// Remove me when deleting `prep2p`
type PREP2PConfig struct {
	ConsensusPort uint32 `json:"consensus_port"`
	DebugPort     uint32 `json:"debug_port"`
}

type P2PConfig struct {
	Protocol   string   `json:"protocol"`
	Address    string   `json:"address"`
	ExternalIp string   `json:"external_ip"`
	Peers      []string `json:"peers"`
}

type ConsensusConfig struct {
	// TODO: This should be set dynamically
	NodeId uint32 `json:"node_id"`
}

type PersistenceConfig struct {
	DataDir string `json:"datadir"`
}

type UtilityConfig struct {
}

func LoadConfig(file string) (c *Config) {
	c = &Config{}

	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatalln("Error opening config file: ", err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("Error reading config file: ", err)
	}
	if err = json.Unmarshal(bytes, c); err != nil {
		log.Fatalln("Error parsing config file: ", err)
	}

	if err := c.validateAndComplete(); err != nil {
		log.Fatalln("Error validating or completing config: ", err)
	}

	if err := c.Consensus.validateAndComplete(); err != nil {
		log.Fatalln("Error validating or completing consensus config: ", err)
	}

	if err := c.P2P.validateAndComplete(); err != nil {
		log.Fatalln("Error validating or completing P2P config: ", err)
	}

	// log.Printf("[DEBUG] ~~~ Finished loading config ~~~\n\t[Config] %#v\n\t[Consensus] %#v\n\t[P2P] %#v\n\t[persistence] %#v\n\t[Utility] %#v\n", c, c.Consensus, c.P2P, c.persistence, c.Utility)

	return
}

func (c *Config) validateAndComplete() error {
	if c.PrivateKey == "" {
		return fmt.Errorf("private key in config file cannot be empty")
	}

	if len(c.Genesis) == 0 {
		return fmt.Errorf("must specify a genesis file")
	}
	c.Genesis = rootify(c.Genesis, c.RootDir)

	return nil
}

func (c *P2PConfig) validateAndComplete() error {
	// if c.ConsensusPort == 0 || c.DebugPort == 0 {
	// 	return fmt.Errorf("ConsensusPort and DebugPort must both be positive integers")
	// }
	return nil
}

func (c *ConsensusConfig) validateAndComplete() error {
	// TODO: c.NodeId should be set dynamically but set via config for testing

	return nil
}

// Helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}