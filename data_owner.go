package feature

// ================================================================
//
// ================================================================
type DataOwnerHttpMethods struct {
	*Feature
}

type DataOwnerEndpoint struct {
	*Endpoint
}

func newDataOwnerEndpoint(e *Endpoint) *DataOwnerEndpoint {
	return &DataOwnerEndpoint{
		Endpoint: e,
	}
}

func (f *Feature) ByAuthorityOfDataOwner() *DataOwnerHttpMethods {
	return &DataOwnerHttpMethods{
		Feature: f,
	}
}

func (m *DataOwnerHttpMethods) GET(relativePath, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "GET", relativePath)
	m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) POST(relativePath, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "POST", relativePath)
	m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) PUT(relativePath, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "PUT", relativePath)
	m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) PATCH(relativePath, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "PATCH", relativePath)
	m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) DELETE(relativePath, identifier string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(identifier, ByAuthorityOfDataOwner, "DELETE", relativePath)
	m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}
