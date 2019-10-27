package server

type torrent struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Magnet string `json:"magnet"`
	URL    string `json:"url"`
}
