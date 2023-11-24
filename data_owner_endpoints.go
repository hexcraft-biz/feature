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

func (m *DataOwnerHttpMethods) GET(relativePath string, scopes []string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(ByAuthorityOfDataOwner, "GET", relativePath, scopes)
	m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) POST(relativePath string, scopes []string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(ByAuthorityOfDataOwner, "POST", relativePath, scopes)
	m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(ByAuthorityOfDataOwner, "PUT", relativePath, scopes)
	m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(ByAuthorityOfDataOwner, "PATCH", relativePath, scopes)
	m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}

func (m *DataOwnerHttpMethods) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) *DataOwnerEndpoint {
	e := m.addEndpoint(ByAuthorityOfDataOwner, "DELETE", relativePath, scopes)
	m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
	return newDataOwnerEndpoint(e)
}
