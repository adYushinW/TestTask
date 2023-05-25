package model

type Cats struct {
	Name            string `json:"name"`
	Color           string `json:"color"`
	Tail_length     int64  `json:"tail_length"`
	Whiskers_length int64  `json:"whiskers_length"`
}

type Cat_colors_info struct {
	Color string `json:"color"`
	Count int64  `json:"count"`
}
