package transcode

type Transcodable interface {
	Flatten() (map[string]interface{}, error)
}
