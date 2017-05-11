package provider

type Provider interface {
	Read() (interface{}, error)
}
