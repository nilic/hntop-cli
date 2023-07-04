package main

import "time"

type Hit struct {
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Author      string    `json:"author"`
	Points      int       `json:"points"`
	NumComments int       `json:"num_comments"`
	ObjectID    string    `json:"objectID"`
}

type Hits struct {
	Hits        []Hit `json:"hits"`
	NbHits      int   `json:"nbHits"`
	Page        int   `json:"page"`
	NbPages     int   `json:"nbPages"`
	HitsPerPage int   `json:"hitsPerPage"`
}
