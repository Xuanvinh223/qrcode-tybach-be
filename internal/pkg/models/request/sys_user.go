package request

type SysUserCreateRequest struct {
	UserName string `form:"username" json:"username" binding:"required"`
	RealName string `form:"realname" json:"realname" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	State    int    `form:"state" json:"state" binding:"required,number,gte=1,lte=2"`
}

type SysUserUpdateRequest struct {
	ID       uint64 `form:"id" json:"id" binding:"required"`
	RealName string `form:"realname" json:"realname" binding:"required"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	State    int    `form:"state" json:"state" binding:"required,number,gte=1,lte=2"`
}

type SysUserDeleteRequest struct {
	ID []uint64 `form:"id" json:"id" binding:"required"`
}

type SysUserLoginByPasswordRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
