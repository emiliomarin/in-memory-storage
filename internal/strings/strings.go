package strings

type GetResponse struct {
	Value     string `json:"value"`
	ExpiresAt string `json:"expires_at,omitempty"`
}
