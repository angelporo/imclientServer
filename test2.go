package main

import (
	"encoding/json"
	"log"
)

type Data struct {
	Page          int
	Pages         int
	PerPage       string
	Total         int
	CountriesList []Country
}
type Country struct {
	Id  string
	Iso string
}

func main() {
	body := []byte(`[
    {
	"page": 1,
	"pages": 6,
	"per_page": "50",
	"total": 256
    },
    [
	{
	    "id": "ABW",
	    "iso2Code": "AW"}]]`)

	items := make([]Data, 10)

	if err := json.Unmarshal(body, &items); err != nil {
		log.Fatalf("error %v", err)
	}
}
