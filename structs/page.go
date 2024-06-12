package structs

type Metadata struct {
	Title          string `json:"title"`
	Keywords       string `json:"keywords"`
	Description    string `json:"description"`
	IsNeedToRender bool   `json:"-"`
}

func (metadata Metadata) SetIsNeedToRender(isNeedToRender bool) Metadata {
	metadata.IsNeedToRender = isNeedToRender

	return metadata
}
