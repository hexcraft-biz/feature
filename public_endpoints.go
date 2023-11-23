package feature

// ================================================================
//
// ================================================================
type PublicHttpMethods struct {
	*Feature
}

type PublicEndpoint struct {
	*Endpoint
}

func newPublicEndpoint(e *Endpoint) *PublicEndpoint {
	return &PublicEndpoint{
		Endpoint: e,
	}
}

func (f *Feature) ByAuthorityOfNone() *PublicHttpMethods {
	return &PublicHttpMethods{
		Feature: f,
	}
}

func (m *PublicHttpMethods) GET(relativePath string, scopes []string, handlers ...HandlerFunc) *PublicEndpoint {
	e := m.addEndpoint(ByAuthorityOfNone, "GET", relativePath, scopes)
	m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
	return newPublicEndpoint(e)
}

func (m *PublicHttpMethods) POST(relativePath string, scopes []string, handlers ...HandlerFunc) *PublicEndpoint {
	e := m.addEndpoint(ByAuthorityOfNone, "POST", relativePath, scopes)
	m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
	return newPublicEndpoint(e)
}

func (m *PublicHttpMethods) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) *PublicEndpoint {
	e := m.addEndpoint(ByAuthorityOfNone, "PUT", relativePath, scopes)
	m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
	return newPublicEndpoint(e)
}

func (m *PublicHttpMethods) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) *PublicEndpoint {
	e := m.addEndpoint(ByAuthorityOfNone, "PATCH", relativePath, scopes)
	m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
	return newPublicEndpoint(e)
}

func (m *PublicHttpMethods) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) *PublicEndpoint {
	e := m.addEndpoint(ByAuthorityOfNone, "DELETE", relativePath, scopes)
	m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
	return newPublicEndpoint(e)
}
