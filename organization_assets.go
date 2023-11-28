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
	e := m.addEndpoint(EnumOwnershipOrganization, "GET", relativePath, scopes, 0)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipOrganization, "POST", relativePath, scopes, 0)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipOrganization, "PUT", relativePath, scopes, 0)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipOrganization, "PATCH", relativePath, scopes, 0)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(EnumOwnershipOrganization, "DELETE", relativePath, scopes, 0)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
