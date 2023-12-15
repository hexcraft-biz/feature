package feature

import (
	"github.com/gin-gonic/gin"
)

// ================================================================
//
// ================================================================
type OrganizationAssets struct {
	*Feature
}

func (m *OrganizationAssets) GET(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipOrganization, "GET", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.GET(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *OrganizationAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipOrganization, "POST", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.POST(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *OrganizationAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipOrganization, "PUT", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.PUT(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *OrganizationAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipOrganization, "PATCH", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(eh, handlers)...)
}

func (m *OrganizationAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	eh := m.addEndpoint(OwnershipOrganization, "DELETE", relativePath, scopes)
	if m.RouterGroup == nil {
		return nil
	}
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(eh, handlers)...)
}
