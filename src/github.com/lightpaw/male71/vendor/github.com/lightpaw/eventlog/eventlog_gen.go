package eventlog

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import "github.com/tinylib/msgp/msgp"

// DecodeMsg implements msgp.Decodable
func (z *data) DecodeMsg(dc *msgp.Reader) (err error) {
	var zjfb uint32
	zjfb, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	if (*z) == nil && zjfb > 0 {
		(*z) = make(data, zjfb)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zjfb > 0 {
		zjfb--
		var zdaf string
		var zpks interface{}
		zdaf, err = dc.ReadString()
		if err != nil {
			return
		}
		zpks, err = dc.ReadIntf()
		if err != nil {
			return
		}
		(*z)[zdaf] = zpks
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z data) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zcxo, zeff := range z {
		err = en.WriteString(zcxo)
		if err != nil {
			return
		}
		err = en.WriteIntf(zeff)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z data) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, uint32(len(z)))
	for zcxo, zeff := range z {
		o = msgp.AppendString(o, zcxo)
		o, err = msgp.AppendIntf(o, zeff)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *data) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zdnj uint32
	zdnj, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	if (*z) == nil && zdnj > 0 {
		(*z) = make(data, zdnj)
	} else if len((*z)) > 0 {
		for key, _ := range *z {
			delete((*z), key)
		}
	}
	for zdnj > 0 {
		var zrsw string
		var zxpk interface{}
		zdnj--
		zrsw, bts, err = msgp.ReadStringBytes(bts)
		if err != nil {
			return
		}
		zxpk, bts, err = msgp.ReadIntfBytes(bts)
		if err != nil {
			return
		}
		(*z)[zrsw] = zxpk
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z data) Msgsize() (s int) {
	s = msgp.MapHeaderSize
	if z != nil {
		for zobc, zsnv := range z {
			_ = zsnv
			s += msgp.StringPrefixSize + len(zobc) + msgp.GuessSize(zsnv)
		}
	}
	return
}
