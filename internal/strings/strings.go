package strings

type GetResponse struct {
	Value     string `json:"value"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int64  `json:"ttl,omitempty"`
}
