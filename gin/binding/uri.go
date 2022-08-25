package binding

type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) BindURI(m map[string][]string, obj any) error {
	return mapURI(obj, m)
}
