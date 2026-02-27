package transcode

type Transcodable interface {
	Flatten() (map[string]any, error)
}
