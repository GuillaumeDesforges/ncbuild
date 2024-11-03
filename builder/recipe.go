package builder

import (
	"crypto/sha256"
	"encoding/hex"
)

type Recipe struct {
	Name       string   `json:""`
	Executable string   `json:""`
	Args       []string `json:""`
}

func (r Recipe) Hash() string {
	summary := r.Name + r.Executable
	hash := sha256.New()
	hash.Write([]byte(summary))
	strHash := hex.EncodeToString(hash.Sum(nil))[:16]
	return strHash
}
