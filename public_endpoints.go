package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PublicHttpMethods struct {
	*Feature
}

func (m *PublicHttpMethods) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfNone, "GET", relativePath, scopes)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicHttpMethods) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfNone, "POST", relativePath, scopes)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicHttpMethods) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfNone, "PUT", relativePath, scopes)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicHttpMethods) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfNone, "PATCH", relativePath, scopes)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicHttpMethods) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(ByAuthorityOfNone, "DELETE", relativePath, scopes)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
