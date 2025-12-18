package metadata

type CategoryResponse struct {
	Data    []CategoryInfo `json:"data"`
	ShowMax int            `json:"show_max"`
}

type CategoryInfo struct {
	ID              int            `json:"id"`
	ParentID        *int           `json:"parent_id,omitempty"`
	Name            string         `json:"name"`
	Type            *string        `json:"type"`
	URLKey          string         `json:"url_key"`
	URLPath         string         `json:"url_path"`
	Level           int            `json:"level"`
	Status          string         `json:"status"`
	IncludeInMenu   any            `json:"include_in_menu"` // be careful, parent is boolean, but children is string :(
	ProductCount    int            `json:"product_count"`
	IsLeaf          bool           `json:"is_leaf"`
	MetaTitle       string         `json:"meta_title"`
	MetaDescription string         `json:"meta_description"`
	MetaKeywords    *string        `json:"meta_keywords"`
	ThumbnailURL    string         `json:"thumbnail_url"`
	FullUrlKey      string         `json:"full_url_key"`
	Children        []CategoryInfo `json:"children"`
}
