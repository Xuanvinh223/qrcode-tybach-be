package types

type SysMenuTreeList struct {
	ID          uint64            `json:"id"`
	MenuName    string            `json:"menuName"`
	MenuPid     string            `json:"menuPid"`
	Url         string            `json:"url"`
	Sort        string            `json:"sort"`
	Description string            `json:"description"`
	State       int               `json:"state"`
	Children    []SysMenuTreeList `json:"children"`
}
