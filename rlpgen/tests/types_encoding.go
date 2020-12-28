package types
	
	import (
		"fmt"

		"github.com/umbracle/fastrlp"
	)

	func (t *Test1) MarshalRLP() []byte {
		return t.MarshalRLPTo(nil)
	}
	
	func (t *Test1) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = t.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (t *Test1) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'A'
vv.Set(ar.NewBytes(t.A[:]))
		
		// Field 'B'
vv.Set(ar.NewBytes(t.B[:]))
		
		// Field 'C'
vv.Set(ar.NewCopyBytes(t.C))
		
		// Field 'D'
vv.Set(ar.NewUint(t.D))
		
		return vv
	}

func (t *Test1) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := t.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (t *Test1) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 4 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'A'
if err = elems[0].GetHash(t.A[:]); err != nil {
				return err
			}
		
		// Field 'B'
if err = elems[1].GetHash(t.B[:]); err != nil {
				return err
			}
		
		// Field 'C'
if t.C, err = elems[2].GetBytes(t.C[:0]); err != nil {
			return err
		}
		
		// Field 'D'
if t.D, err = elems[3].GetUint64(); err != nil {
			return err
		}
		
		return nil
	}
func (h *Header) MarshalRLP() []byte {
		return h.MarshalRLPTo(nil)
	}
	
	func (h *Header) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = h.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (h *Header) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'ParentHash'
vv.Set(ar.NewBytes(h.ParentHash[:]))
		
		// Field 'Sha3Uncles'
vv.Set(ar.NewBytes(h.Sha3Uncles[:]))
		
		// Field 'Miner'
vv.Set(ar.NewBytes(h.Miner[:]))
		
		// Field 'StateRoot'
vv.Set(ar.NewBytes(h.StateRoot[:]))
		
		// Field 'TxRoot'
vv.Set(ar.NewBytes(h.TxRoot[:]))
		
		// Field 'ReceiptsRoot'
vv.Set(ar.NewBytes(h.ReceiptsRoot[:]))
		
		// Field 'LogsBloom'
vv.Set(ar.NewBytes(h.LogsBloom[:]))
		
		// Field 'Difficulty'
vv.Set(ar.NewUint(h.Difficulty))
		
		// Field 'Number'
vv.Set(ar.NewUint(h.Number))
		
		// Field 'GasLimit'
vv.Set(ar.NewUint(h.GasLimit))
		
		// Field 'GasUsed'
vv.Set(ar.NewUint(h.GasUsed))
		
		// Field 'Timestamp'
vv.Set(ar.NewUint(h.Timestamp))
		
		// Field 'ExtraData'
vv.Set(ar.NewCopyBytes(h.ExtraData))
		
		// Field 'MixHash'
vv.Set(ar.NewBytes(h.MixHash[:]))
		
		// Field 'Nonce'
vv.Set(ar.NewBytes(h.Nonce[:]))
		
		return vv
	}

func (h *Header) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := h.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (h *Header) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 15 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'ParentHash'
if err = elems[0].GetHash(h.ParentHash[:]); err != nil {
				return err
			}
		
		// Field 'Sha3Uncles'
if err = elems[1].GetHash(h.Sha3Uncles[:]); err != nil {
				return err
			}
		
		// Field 'Miner'
if err = elems[2].GetAddr(h.Miner[:]); err != nil {
				return err
			}
		
		// Field 'StateRoot'
if err = elems[3].GetHash(h.StateRoot[:]); err != nil {
				return err
			}
		
		// Field 'TxRoot'
if err = elems[4].GetHash(h.TxRoot[:]); err != nil {
				return err
			}
		
		// Field 'ReceiptsRoot'
if err = elems[5].GetHash(h.ReceiptsRoot[:]); err != nil {
				return err
			}
		
		// Field 'LogsBloom'
if _, err = elems[6].GetBytes(h.LogsBloom[:0], 256); err != nil {
				return err
			}
		
		// Field 'Difficulty'
if h.Difficulty, err = elems[7].GetUint64(); err != nil {
			return err
		}
		
		// Field 'Number'
if h.Number, err = elems[8].GetUint64(); err != nil {
			return err
		}
		
		// Field 'GasLimit'
if h.GasLimit, err = elems[9].GetUint64(); err != nil {
			return err
		}
		
		// Field 'GasUsed'
if h.GasUsed, err = elems[10].GetUint64(); err != nil {
			return err
		}
		
		// Field 'Timestamp'
if h.Timestamp, err = elems[11].GetUint64(); err != nil {
			return err
		}
		
		// Field 'ExtraData'
if h.ExtraData, err = elems[12].GetBytes(h.ExtraData[:0]); err != nil {
			return err
		}
		
		// Field 'MixHash'
if err = elems[13].GetHash(h.MixHash[:]); err != nil {
				return err
			}
		
		// Field 'Nonce'
if _, err = elems[14].GetBytes(h.Nonce[:0], 8); err != nil {
				return err
			}
		
		return nil
	}
func (t *Transaction) MarshalRLP() []byte {
		return t.MarshalRLPTo(nil)
	}
	
	func (t *Transaction) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = t.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (t *Transaction) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'Nonce'
vv.Set(ar.NewUint(t.Nonce))
		
		// Field 'GasPrice'
vv.Set(ar.NewCopyBytes(t.GasPrice))
		
		// Field 'Gas'
vv.Set(ar.NewUint(t.Gas))
		
		// Field 'To'
if t.To == nil {
                vv.Set(ar.NewNull())
            } else {
                vv.Set(ar.NewBytes(t.To[:]))
            }
		
		// Field 'Value'
vv.Set(ar.NewCopyBytes(t.Value))
		
		// Field 'Input'
vv.Set(ar.NewCopyBytes(t.Input))
		
		// Field 'V'
vv.Set(ar.NewUint(uint64(t.V)))
		
		// Field 'R'
vv.Set(ar.NewCopyBytes(t.R))
		
		// Field 'S'
vv.Set(ar.NewCopyBytes(t.S))
		
		return vv
	}

func (t *Transaction) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := t.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (t *Transaction) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 9 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'Nonce'
if t.Nonce, err = elems[0].GetUint64(); err != nil {
			return err
		}
		
		// Field 'GasPrice'
if t.GasPrice, err = elems[1].GetBytes(t.GasPrice[:0]); err != nil {
			return err
		}
		
		// Field 'Gas'
if t.Gas, err = elems[2].GetUint64(); err != nil {
			return err
		}
		
		// Field 'To'
if err = elems[3].GetAddr(t.To[:]); err != nil {
				return err
			}
		
		// Field 'Value'
if t.Value, err = elems[4].GetBytes(t.Value[:0]); err != nil {
			return err
		}
		
		// Field 'Input'
if t.Input, err = elems[5].GetBytes(t.Input[:0]); err != nil {
			return err
		}
		
		// Field 'V'
if t.V, err = elems[6].GetByte(); err != nil {
			return err
		}
		
		// Field 'R'
if t.R, err = elems[7].GetBytes(t.R[:0]); err != nil {
			return err
		}
		
		// Field 'S'
if t.S, err = elems[8].GetBytes(t.S[:0]); err != nil {
			return err
		}
		
		return nil
	}
func (b *Body) MarshalRLP() []byte {
		return b.MarshalRLPTo(nil)
	}
	
	func (b *Body) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = b.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (b *Body) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'Transactions'
{
			if len(b.Transactions) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range b.Transactions {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}
		
		// Field 'Uncles'
{
			if len(b.Uncles) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range b.Uncles {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}
		
		return vv
	}

func (b *Body) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := b.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (b *Body) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 2 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'Transactions'
{
			subElems, err := elems[0].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := &Transaction{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				b.Transactions = append(b.Transactions, bb)
			}
		}
		
		// Field 'Uncles'
{
			subElems, err := elems[1].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := &Header{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				b.Uncles = append(b.Uncles, bb)
			}
		}
		
		return nil
	}
func (b *Block) MarshalRLP() []byte {
		return b.MarshalRLPTo(nil)
	}
	
	func (b *Block) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = b.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (b *Block) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'Header'
vv.Set(b.Header.MarshalRLPWith(ar))
		
		// Field 'Transactions'
{
			if len(b.Transactions) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range b.Transactions {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}
		
		// Field 'Uncles'
{
			if len(b.Uncles) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range b.Uncles {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}
		
		return vv
	}

func (b *Block) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := b.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (b *Block) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 3 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'Header'
{
			b.Header = &Header{}
			if err := b.Header.UnmarshalRLPFrom(elems[0]); err != nil {
				return err
			}
		}
		
		// Field 'Transactions'
{
			subElems, err := elems[1].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := &Transaction{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				b.Transactions = append(b.Transactions, bb)
			}
		}
		
		// Field 'Uncles'
{
			subElems, err := elems[2].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := &Header{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				b.Uncles = append(b.Uncles, bb)
			}
		}
		
		return nil
	}
func (r *Receipt) MarshalRLP() []byte {
		return r.MarshalRLPTo(nil)
	}
	
	func (r *Receipt) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = r.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (r *Receipt) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'Root'
vv.Set(ar.NewCopyBytes(r.Root))
		
		// Field 'CumulativeGasUsed'
vv.Set(ar.NewUint(r.CumulativeGasUsed))
		
		// Field 'Logs'
{
			if len(r.Logs) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range r.Logs {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}
		
		return vv
	}

func (r *Receipt) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := r.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (r *Receipt) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 3 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'Root'
if r.Root, err = elems[0].GetBytes(r.Root[:0]); err != nil {
			return err
		}
		
		// Field 'CumulativeGasUsed'
if r.CumulativeGasUsed, err = elems[1].GetUint64(); err != nil {
			return err
		}
		
		// Field 'Logs'
{
			subElems, err := elems[2].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := &Log{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				r.Logs = append(r.Logs, bb)
			}
		}
		
		return nil
	}
func (l *Log) MarshalRLP() []byte {
		return l.MarshalRLPTo(nil)
	}
	
	func (l *Log) MarshalRLPTo(dst []byte) []byte {
		ar := fastrlp.DefaultArenaPool.Get()
		dst = l.MarshalRLPWith(ar).MarshalTo(dst)
		fastrlp.DefaultArenaPool.Put(ar)
		return dst
	}

	func (l *Log) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
		vv := ar.NewArray()
		
		// Field 'Address'
vv.Set(ar.NewBytes(l.Address[:]))
		
		// Field 'Topics'
{
			if len(l.Topics) == 0 {
				vv.Set(ar.NewNullArray())
			} else {
				v0 := ar.NewArray()
				for _, item := range l.Topics {
					v0.Set(item.MarshalRLPWith(ar))
				}
				vv.Set(v0)
			}
		}
		
		// Field 'Data'
vv.Set(ar.NewCopyBytes(l.Data))
		
		return vv
	}

func (l *Log) UnmarshalRLP(buf []byte) error {
		pr := fastrlp.DefaultParserPool.Get()
		defer fastrlp.DefaultParserPool.Put(pr)

		vv, err := pr.Parse(buf)
		if err != nil {
			return err
		}
		if err := l.UnmarshalRLPFrom(vv); err != nil {
			return err
		}
		return nil
	}

	func (l *Log) UnmarshalRLPFrom(v *fastrlp.Value) error {
		elems, err := v.GetElems()
		if err != nil {
			return err
		}
		if num := len(elems); num != 3 {
			return fmt.Errorf("not enough elements to decode transaction, expected 9 but found %d", num)
		}
		
		// Field 'Address'
if err = elems[0].GetAddr(l.Address[:]); err != nil {
				return err
			}
		
		// Field 'Topics'
{
			subElems, err := elems[1].GetElems()
			if err != nil {
				return err
			}
			for _, elem := range subElems {
				bb := &{}
				if err := bb.UnmarshalRLPFrom(elem); err != nil {
					return err
				}
				l.Topics = append(l.Topics, bb)
			}
		}
		
		// Field 'Data'
if l.Data, err = elems[2].GetBytes(l.Data[:0]); err != nil {
			return err
		}
		
		return nil
	}
	