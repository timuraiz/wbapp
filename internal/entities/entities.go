package entities

type Category struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type CategoriesResponse struct {
	Data struct {
		Categories []Category `json:"categories"`
	} `json:"data"`
}
