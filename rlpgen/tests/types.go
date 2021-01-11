package types

import "github.com/umbracle/fastrlp/rlpgen/tests/external"

type Hash [32]byte

type Address [20]byte

const (
	// sszgen cannot decode this num so the user
	// has to supply the length for Bloom as "--sizes Bloom=256"
	num = 256
)

type Bloom [num]byte

type Test1 struct {
	A external.Fixed
	B [32]byte
	C []byte
	D uint64
}

func a() {

}

type Header struct {
	ParentHash   Hash
	Sha3Uncles   Hash
	Miner        Address
	StateRoot    Hash
	TxRoot       Hash
	ReceiptsRoot Hash
	LogsBloom    Bloom
	Difficulty   uint64
	Number       uint64
	GasLimit     uint64
	GasUsed      uint64
	Timestamp    uint64
	ExtraData    []byte
	MixHash      Hash
	Nonce        Nonce

	Hash Hash `rlp:"-"` // this field has to be avoided
}

type Transaction struct {
	Nonce    uint64
	GasPrice []byte
	Gas      uint64
	To       *Address
	Value    []byte
	Input    []byte

	V byte
	R []byte
	S []byte
}

type Body struct {
	Transactions []*Transaction
	Uncles       []*Header
}

type Block struct {
	Header       *Header
	Transactions []*Transaction
	Uncles       []*Header
}

type Receipt struct {
	Root              []byte
	CumulativeGasUsed uint64
	Logs              []*Log
}

type Log struct {
	Address Address
	Topics  []Hash
	Data    []byte
}
