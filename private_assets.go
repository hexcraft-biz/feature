package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PrivateAssets struct {
	*Feature
}

func (m *PrivateAssets) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPrivate, "GET", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.GET(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PrivateAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPrivate, "POST", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.POST(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PrivateAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPrivate, "PUT", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.PUT(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PrivateAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPrivate, "PATCH", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PrivateAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPrivate, "DELETE", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(eh, handlers)...)
}
