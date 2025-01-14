// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package test_helpers

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/ens_transformers/test_config"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/contract"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

func SetupDBandBC() (*postgres.DB, core.BlockChain) {
	cli := test_config.InfuraClient
	infuraIPC := cli.IPCPath
	rawRpcClient, err := rpc.Dial(infuraIPC)
	Expect(err).NotTo(HaveOccurred())
	rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
	ethClient := ethclient.NewClient(rawRpcClient)
	blockChainClient := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, blockChain.Node())
	Expect(err).NotTo(HaveOccurred())

	return db, blockChain
}

func SetupENSRepo(start int64) (*postgres.DB, *contract.Contract) {
	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, core.Node{})
	Expect(err).NotTo(HaveOccurred())

	info := SetupENSRegistryContract(start)

	return db, info
}

func TearDown(db *postgres.DB) {
	tx, err := db.Beginx()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP TABLE checked_headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE checked_headers (id SERIAL PRIMARY KEY, header_id INTEGER UNIQUE NOT NULL REFERENCES headers (id) ON DELETE CASCADE);`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS auction_started_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS bid_revealed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS hash_invalidated_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS hash_registered_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS hash_released_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS new_bid_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS new_owner_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS new_resolver_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS new_ttl_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS ens_transfer_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS abi_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS addr_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS content_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS contenthash_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS multihash_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS name_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS pubkey_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS text_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())
	_, err = tx.Exec(`ALTER TABLE checked_headers ADD column IF NOT EXISTS name_changed_checked INTEGER NOT NULL DEFAULT 0`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM ens.domain_records`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())
}

func SetupENSRegistryContract(start int64) *contract.Contract {
	p := mocks.NewParser(constants.ENSAbiString)
	err := p.Parse()
	Expect(err).ToNot(HaveOccurred())

	return contract.Contract{
		Name:          "ENS-Registry",
		Network:       "",
		Address:       constants.EnsContractAddress,
		Abi:           p.Abi(),
		ParsedAbi:     p.ParsedAbi(),
		StartingBlock: start,
		Events:        p.GetEvents([]string{}),
		Methods:       nil,
		FilterArgs:    map[string]bool{},
		MethodArgs:    map[string]bool{},
	}.Init()
}
