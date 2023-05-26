package model

type Cats struct {
	Name            string `json:"name"`
	Color           string `json:"color"`
	Tail_length     uint64 `json:"tail_length"`
	Whiskers_length uint64 `json:"whiskers_length"`
}

type Cat_colors_info struct {
	Color string `json:"color"`
	Count uint64 `json:"count"`
}

type Cats_stat struct {
	Tail_length_mean       float64 `json:"tail_length_mean"`
	Tail_length_median     float64 `json:"tail_length_median"`
	Tail_length_mode       string  `json:"tail_length_mode"`
	Whiskers_length_mean   float64 `json:"whiskers_length_mean"`
	Whiskers_length_median float64 `json:"whiskers_length_median"`
	Whiskers_length_mode   string  `json:"whiskers_length_mode"`
}
