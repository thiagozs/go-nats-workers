package models

type Message struct {
	Order       string `cbor:"1,keyasint,omitempty"`
	Matches     string `cbor:"2,keyasint,omitempty"`
	Operations  string `cbor:"3,keyasint,omitempty"`
	Commissions string `cbor:"4,keyasint,omitempty"`
}
