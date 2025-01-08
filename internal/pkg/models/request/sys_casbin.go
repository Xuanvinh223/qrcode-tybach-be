package request

type RouteRequest struct {
	RType  string `form:"rtype" json:"rtype" binding:"required"`   // role url type (/api/v1, /api/v1/users)
	RoleId int    `form:"roleId" json:"roleId" binding:"required"` // role id
}

type SysCasbinCreateRequest struct {
	RoleId uint64 `form:"roleId" json:"roleId" binding:"required"`
	Url    string `form:"url" json:"url" binding:"required"`
	Method string `form:"method" json:"method" binding:"required"`
}

type SysCasbinUpdateRequest struct {
	ID     uint64 `form:"id" json:"id" binding:"required"`
	RoleId int    `form:"roleId" json:"roleId" binding:"required"`
	Url    string `form:"url" json:"url" binding:"required"`
	Method string `form:"method" json:"method" binding:"required"`
}

type SysCasbinDeleteRequest struct {
	ID []uint64 `form:"id" json:"id" binding:"required"`
}
