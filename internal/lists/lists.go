package lists

type GetResponse[T any] struct {
	List      []T    `json:"list"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type SetRequest[T any] struct {
	Key  string `json:"key"`
	List []T    `json:"list"`
	TTL  int64  `json:"ttl,omitempty"`
}

type PopResponse[T any] struct {
	Value T `json:"value"`
}

type UpdateRequest[T any] struct {
	Key  string `json:"key"`
	List []T    `json:"list"`
}
