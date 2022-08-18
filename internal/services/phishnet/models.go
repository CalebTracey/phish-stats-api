package phishnet

type ShowResponse struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
	Data         []struct {
		Showid              string `json:"showid"`
		Showdate            string `json:"showdate"`
		Permalink           string `json:"permalink"`
		Showyear            string `json:"showyear"`
		Uniqueid            string `json:"uniqueid"`
		Meta                string `json:"meta"`
		Reviews             string `json:"reviews"`
		Exclude             string `json:"exclude"`
		Setlistnotes        string `json:"setlistnotes"`
		Soundcheck          string `json:"soundcheck"`
		Songid              string `json:"songid"`
		Position            string `json:"position"`
		Transition          string `json:"transition"`
		Footnote            string `json:"footnote"`
		Set                 string `json:"set"`
		Isjam               string `json:"isjam"`
		Isreprise           string `json:"isreprise"`
		Isjamchart          string `json:"isjamchart"`
		JamchartDescription string `json:"jamchart_description"`
		Tracktime           string `json:"tracktime"`
		Gap                 string `json:"gap"`
		Tourid              string `json:"tourid"`
		Tourname            string `json:"tourname"`
		Tourwhen            string `json:"tourwhen"`
		Song                string `json:"song"`
		Nickname            string `json:"nickname"`
		Slug                string `json:"slug"`
		IsOriginal          string `json:"is_original"`
		Venueid             string `json:"venueid"`
		Venue               string `json:"venue"`
		City                string `json:"city"`
		State               string `json:"state"`
		Country             string `json:"country"`
		TransMark           string `json:"trans_mark"`
		Artistid            string `json:"artistid"`
		ArtistSlug          string `json:"artist_slug"`
		ArtistName          string `json:"artist_name"`
	} `json:"data"`
}
