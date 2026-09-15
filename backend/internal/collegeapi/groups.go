package collegeapi

type Group struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type groupsResponse struct {
	Results []Group `json:"results"`
}