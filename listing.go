package geddit

// ListingMeta contains the metadata passed back from the json request
type ListingMeta struct {
	After   string `json:"after,omitempty"`
	Before  string `json:"before,omitempty"`
	Dist    int    `json:"dist,omitempty"`
	ModHash string `json:"modHash,omitempty"`
}

type ListingResp struct {
	Data struct {
		Children []struct {
			Data *Submission
		}
		ListingMeta
	}
}
