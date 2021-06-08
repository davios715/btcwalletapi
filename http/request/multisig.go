package request

type MultiSig struct {
	N          int      `json:"n"`
	M          int      `json:"m"`
	PublicKeys []string `json:"public_keys"`
}
