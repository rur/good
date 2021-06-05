package generate

// Handler is data for a handler function which should be created
type Handler struct {
	Ref        string
	Extends    string
	Method     string
	Doc        string
	Identifier string
	Blocks     []HandleBlock
}

// HandleBlock is the details of sub-views which should
// be delegated to in the handler
type HandleBlock struct {
	FieldName string
	Name      string
}
