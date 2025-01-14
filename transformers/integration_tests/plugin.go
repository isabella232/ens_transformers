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

package integration_tests

import (
	"plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"

	"github.com/vulcanize/ens_transformers/test_config"
	"github.com/vulcanize/ens_transformers/transformers/domain_records/test_helpers"
)

var eventConfig = config.Plugin{
	Home: "github.com/vulcanize/ens_transformers",
	Transformers: map[string]config.Transformer{
		"newOwner": {
			Path:           "transformers/registry/new_owner/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/vulcanize/ens_transformers",
		},
		"addrChanged": {
			Path:           "transformers/resolver/addr_changed/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/vulcanize/ens_transformers",
		},
		"newBid": {
			Path:           "transformers/registar/new_bid/initializer",
			Type:           config.EthEvent,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/vulcanize/ens_transformers",
		},
	},
	FileName: "testEventTransformerSet",
	FilePath: "$GOPATH/src/github.com/vulcanize/ens_transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var dbConfig = config.Database{
	Hostname: "localhost",
	Port:     5432,
	Name:     "vulcanize_private",
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer)
}

var _ = Describe("Plugin test", func() {
	var g p2.Generator
	var goPath, soPath string
	var err error
	var bc core.BlockChain
	var db *postgres.DB
	var hr repositories.HeaderRepository
	var headerID int64
	viper.SetConfigName("composeAndExecuteEventTransformers")
	viper.AddConfigPath("$GOPATH/src/github.com/vulcanize/ens_transformers/environments/")

	Describe("Event Transformers only", func() {
		BeforeEach(func() {
			goPath, soPath, err = eventConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(eventConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})
		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified  TransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				initializers, store, _ := exporter.Export()
				Expect(len(initializers)).To(Equal(3))
				Expect(len(store)).To(Equal(0))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of TransformerInitializers that we can execute over", func() {
				db, bc = test_helpers.SetupDBandBC()
				defer test_config.CleanTestDB(db)

				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(7483567)
				Expect(err).ToNot(HaveOccurred())
				headerID, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				initializers, _, _ := exporter.Export()

				w := watcher.NewEventWatcher(db, bc)
				w.AddTransformers(initializers)
				err = w.Execute(constants.HeaderMissing)
				Expect(err).ToNot(HaveOccurred())

				type model struct {
					Node             string
					Label            string
					Owner            string
					Subnode          string
					LogIndex         uint   `db:"log_idx"`
					TransactionIndex uint   `db:"tx_idx"`
					Raw              []byte `db:"raw_log"`
					Id               int64  `db:"id"`
					HeaderId         int64  `db:"header_id"`
				}

				returned := model{}

				err = db.Get(&returned, `SELECT * FROM ens.new_owner WHERE header_id = $1`, headerID)
				Expect(err).ToNot(HaveOccurred())
				Expect(returned.Node).To(Equal("0x79aa1ec377dbfd1edf87d526cd5c116ac6ec4444e23da2a8b8ae0e9db9f46ec9"))
				Expect(returned.Label).To(Equal("0xf06ea683845d8994338d98e7296850612317647cdb1db3d1051f24bafdefe9f7"))
				Expect(returned.Owner).To(Equal("0xA964ed4077aD3ba1946D118ce90544657bB4003B"))
				Expect(returned.Subnode).To(Equal("0xbb87bd9021ba9da3248899e6fdd901a68efb0e15ac691ac9ce5cc88ebcb306de"))
				Expect(returned.TransactionIndex).To(Equal(uint(50)))
				Expect(returned.LogIndex).To(Equal(uint(78)))
			})
		})
	})
})
