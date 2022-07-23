package utility_module

import (
	"fmt"
	"log"
	"math/big"
	"testing"
	"os/signal"
	"os"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/pokt-network/pocket/persistence"
	"github.com/stretchr/testify/require"

	"github.com/pokt-network/pocket/shared/config"
	"github.com/pokt-network/pocket/shared/types"
	typesGenesis "github.com/pokt-network/pocket/shared/types/genesis"
	"github.com/pokt-network/pocket/utility"
)

const (
	user             = "postgres"
	password         = "secret"
	db               = "postgres"
	sql_schema       = "test_schema"
	localhost        = "0.0.0.0"
	port             = "5432"
	dialect          = "postgres"
	connStringFormat = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

var (
	defaultTestingChains          = []string{"0001"}
	defaultTestingChainsEdited    = []string{"0002"}
	defaultServiceUrl             = "https://foo.bar"
	defaultServiceUrlEdited       = "https://bar.foo"
	defaultServiceNodesPerSession = 24
	zeroAmount                    = big.NewInt(0)
	zeroAmountString              = types.BigIntToString(zeroAmount)
	defaultAmount                 = big.NewInt(1000000000000000)
	defaultSendAmount             = big.NewInt(10000)
	defaultAmountString           = types.BigIntToString(defaultAmount)
	defaultNonceString            = types.BigIntToString(defaultAmount)
	defaultSendAmountString       = types.BigIntToString(defaultSendAmount)
)

func NewTestingMempool(_ *testing.T) types.Mempool {
	return types.NewMempool(1000000, 1000)
}

func NewTestingUtilityContext(t *testing.T, height int64) utility.UtilityContext {
	mempool := NewTestingMempool(t)
	cfg := &config.Config{
		Genesis: genesisJson(), 
		Persistence: &config.PersistenceConfig{"", fmt.Sprintf(connStringFormat, user, password, port, db), ""},
	}
	_ = typesGenesis.GetNodeState(cfg)
	persistenceModule, err := persistence.NewPersistenceModule(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := persistenceModule.Start(); err != nil {
		t.Fatal(err)
	}
	persistenceContext, err := persistenceModule.NewContext(height)
	require.NoError(t, err)
	return utility.UtilityContext{
		LatestHeight: height,
		Mempool:      mempool,
		Context: &utility.Context{
			PersistenceContext: persistenceContext,
			SavePointsM:        make(map[string]struct{}),
			SavePoints:         make([][]byte, 0),
		},
	}
}

func genesisJson() string {
	return fmt.Sprintf(`{
		"genesis_state_configs": {
			"num_validators": 5,
			"num_applications": 1,
			"num_fisherman": 1,
			"num_servicers": 5,
			"keys_seed_start": %d
		},
		"genesis_time": "2022-01-19T00:00:00.000000Z",
		"app_hash": "genesis_block_or_state_hash"
	}`, 42)
}

var PostgresDB *persistence.PostgresDB

// TODO(team): make these tests thread safe
func init() {
	PostgresDB = new(persistence.PostgresDB)
}

// See https://github.com/ory/dockertest as reference for the template of this code
func TestMain(m *testing.M) {
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{port},
		PortBindings: map[docker.Port][]docker.PortBinding{
			port: {
				{HostIP: localhost, HostPort: port},
			},
		},
	}
	connString := fmt.Sprintf(connStringFormat, user, password, port, db)

	defer func() {
		ctx, _ := PostgresDB.GetContext()
		PostgresDB.Conn.Close(ctx)
		ctx.Done()
	}()

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("***Make sure your docker daemon is running!!*** Could not start resource: %s\n", err.Error())
	}

	// DOCUMENT: Why do we not call `syscall.SIGTERM` here
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		for sig := range c {
			log.Printf("exit signal %d received\n", sig)
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("could not purge resource: %s", err)
			}
		}
	}()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err = pool.Retry(func() error {
		conn, err := persistence.ConnectAndInitializeDatabase(connString, sql_schema)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		PostgresDB.Conn = conn
		return nil
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}
	code := m.Run()

	// You can't defer this because `os.Exit`` doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Exit(code)
}
