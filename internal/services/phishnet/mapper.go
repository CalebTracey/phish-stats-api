package phishnet

import "github.com/calebtracey/phish-stats-api/internal/models"

//go:generate mockgen -destination=mockMapper.go -package=phishnet . MapperI
type MapperI interface {
	PhishNetResponseToShowResponse(pnResponse PNShowResponse) (response models.ShowResponse)
}

type Mapper struct{}

func (m Mapper) PhishNetResponseToShowResponse(pnResponse PNShowResponse) (response models.ShowResponse) {
	var show models.Show

	if len(pnResponse.Data) > 0 {
		show.Venue = pnResponse.Data[0].Venue
		show.Date = pnResponse.Data[0].Showdate
		for _, data := range pnResponse.Data {
			show.Songs = append(show.Songs, models.Song{
				SongID:    data.Songid,
				Title:     data.Song,
				TrackTime: data.Tracktime,
				Gap:       data.Gap,
			})
		}
	}
	response.Show = show

	return response
}
