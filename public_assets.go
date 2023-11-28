package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PublicAssets struct {
	*Feature
}

func (m *PublicAssets) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPublic, "GET", relativePath, scopes, 0)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPublic, "POST", relativePath, scopes, 0)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPublic, "PUT", relativePath, scopes, 0)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPublic, "PATCH", relativePath, scopes, 0)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PublicAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPublic, "DELETE", relativePath, scopes, 0)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
