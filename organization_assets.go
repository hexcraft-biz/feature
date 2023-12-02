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
	e := m.addEndpoint(OwnershipOrganization, "GET", relativePath, scopes)
	return m.RouterGroup.GET(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) POST(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipOrganization, "POST", relativePath, scopes)
	return m.RouterGroup.POST(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) PUT(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipOrganization, "PUT", relativePath, scopes)
	return m.RouterGroup.PUT(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) PATCH(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipOrganization, "PATCH", relativePath, scopes)
	return m.RouterGroup.PATCH(relativePath, handlerFuncs(e, handlers)...)
}

func (m *OrganizationAssets) DELETE(relativePath string, scopes []string, handlers ...HandlerFunc) gin.IRoutes {
	e := m.addEndpoint(OwnershipOrganization, "DELETE", relativePath, scopes)
	return m.RouterGroup.DELETE(relativePath, handlerFuncs(e, handlers)...)
}
