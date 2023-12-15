package feature

import "github.com/gin-gonic/gin"

// ================================================================
//
// ================================================================
type PublicAssets struct {
	*Feature
}

func (m *PublicAssets) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPublic, "GET", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.GET(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PublicAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPublic, "POST", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.POST(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PublicAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPublic, "PUT", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.PUT(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PublicAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPublic, "PATCH", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *PublicAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipPublic, "DELETE", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(eh, handlers)...)
}
