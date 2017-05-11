package provider

type Provider interface {
	Read() (map[string]interface{}, error)
}
