package file

type FileProvider struct{}

func NewProvider(path string) *FileProvider {
	return &FileProvider{}
}

func (fp *FileProvider) Read() (interface{}, error) {
	return nil, nil
}
