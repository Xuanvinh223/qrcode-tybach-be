package types

type SysRoleInfo struct {
	ID           uint64   `json:"id"`
	RoleName     string   `json:"roleName"`
	Description  string   `json:"description"`
	RoleUserList []uint64 `json:"roleUserList"`
	RoleMenuList []uint64 `json:"roleMenuList"`
}
