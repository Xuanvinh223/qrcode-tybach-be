package types

type RouteInfo struct {
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	NewPath    string   `json:"new_path"`
	MethodList []string `json:"method_list"`
}
