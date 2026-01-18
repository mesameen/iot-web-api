package model

type Entities struct {
	TenantGroups []TenantGroup `json:"tenant_groups"`
	Tenants      []Tenant      `json:"tenants"`
	Parsers      []Parser      `json:"parsers"`
}

type TenantGroup struct {
	Name          string `json:"name"`
	TenantGroupID string `json:"tenant_group_id"`
}

type Tenant struct {
	Name     string `json:"name"`
	TenantID string `json:"tenant_id"`
}

type Parser struct {
	Name     string `json:"name"`
	ParserID int    `json:"parser_id"`
}
