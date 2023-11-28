package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PrivateAssets struct {
	*Feature
}

func (m *PrivateAssets) GET(relativePath string, ownerParamIndex int, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPrivate, "GET", relativePath, scopes, ownerParamIndex)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) POST(relativePath string, ownerParamIndex int, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPrivate, "POST", relativePath, scopes, ownerParamIndex)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) PUT(relativePath string, ownerParamIndex int, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPrivate, "PUT", relativePath, scopes, ownerParamIndex)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) PATCH(relativePath string, ownerParamIndex int, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPrivate, "PATCH", relativePath, scopes, ownerParamIndex)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *PrivateAssets) DELETE(relativePath string, ownerParamIndex int, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipPrivate, "DELETE", relativePath, scopes, ownerParamIndex)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
