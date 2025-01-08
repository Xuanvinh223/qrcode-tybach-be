package request

type SysRoleCreateRequest struct {
	RoleName     string   `form:"rolename" json:"rolename" binding:"required"`
	Description  string   `form:"description" json:"description"`
	RoleUserList []uint64 `form:"roleUserList" json:"roleUserList"`
	RoleMenuList []uint64 `form:"roleMenuList" json:"roleMenuList"`
}

type SysRoleUpdateRequest struct {
	ID           uint64   `form:"id" json:"id" binding:"required"`
	RoleName     string   `form:"rolename" json:"rolename" binding:"required"`
	Description  string   `form:"description" json:"description"`
	RoleUserList []uint64 `form:"roleUserList" json:"roleUserList"`
	RoleMenuList []uint64 `form:"roleMenuList" json:"roleMenuList"`
}

type SysRoleDeleteRequest struct {
	ID []uint64 `form:"id" json:"id" binding:"required"`
}
