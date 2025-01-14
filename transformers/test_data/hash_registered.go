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

package test_data

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/ens_transformers/transformers/registar/hash_registered"
)

const (
	TemporaryHashRegisteredBlockNumber = int64(26)
	TemporaryHashRegisteredData        = "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002"
	TemporaryHashRegisteredTransaction = "0x5c698f13940a2153440c6d19660878bc90219d9298fdcf37365aa8d88d40fc42"
)

var (
	hashRegisteredRawJson, _ = json.Marshal(EthHashRegisteredLog)
)

var EthHashRegisteredLog = types.Log{
	Address: common.HexToAddress(RegistarAddress),
	Topics: []common.Hash{
		common.HexToHash("0x99b5620489b6ef926d4518936cfec15d305452712b88bd59da2d9c10fb0953e8"),
		common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000000000d8b4147eda80fec7122ae16da2479cbd7ffb"),
	},
	Data:        hexutil.MustDecode(TemporaryHashRegisteredData),
	BlockNumber: uint64(TemporaryHashRegisteredBlockNumber),
	TxHash:      common.HexToHash(TemporaryHashRegisteredTransaction),
	TxIndex:     111,
	BlockHash:   fakes.FakeHash,
	Index:       7,
	Removed:     false,
}

var HashRegisteredEntity = hash_registered.HashRegisteredEntity{
	Hash:             node,
	Owner:            owner,
	Value:            value,
	RegistrationDate: registrationDate,
	LogIndex:         EthHashRegisteredLog.Index,
	TransactionIndex: EthHashRegisteredLog.TxIndex,
	Raw:              EthHashRegisteredLog,
}

var HashRegisteredModel = hash_registered.HashRegisteredModel{
	Hash:             node.Hex(),
	Owner:            owner.Hex(),
	Value:            value.String(),
	RegistrationDate: registrationDate.String(),
	LogIndex:         EthHashRegisteredLog.Index,
	TransactionIndex: EthHashRegisteredLog.TxIndex,
	Raw:              hashRegisteredRawJson,
}
