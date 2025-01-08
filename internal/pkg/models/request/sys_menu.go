package request

type SysMenuCreateRequest struct {
	MenuName    string `form:"menuName" json:"menuName" binding:"required"`
	MenuPid     string `form:"menuPid" json:"menuPid" binding:"required"`
	Url         string `form:"url" json:"url"`
	Sort        string `form:"sort" json:"sort"`
	Description string `form:"description" json:"description"`
	State       int    `form:"state" json:"state" binding:"required,number,gte=1,lte=2"`
}

type SysMenuUpdateRequest struct {
	ID          uint64 `form:"id" json:"id" binding:"required"`
	MenuName    string `form:"menuName" json:"menuName" binding:"required"`
	MenuPid     string `form:"menuPid" json:"menuPid" binding:"required"`
	Url         string `form:"url" json:"url"`
	Sort        string `form:"sort" json:"sort"`
	Description string `form:"description" json:"description"`
	State       int    `form:"state" json:"state" binding:"required,number,gte=1,lte=2"`
}

type SysMenuDeleteRequest struct {
	ID []uint64 `form:"id" json:"id" binding:"required"`
}
