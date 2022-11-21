package phishnet

import "github.com/calebtracey/phish-stats-api/internal/models"

//go:generate mockgen -destination=mockMapper.go -package=phishnet . MapperI
type MapperI interface {
	PhishNetResponseToShowResponse(pnResponse PNShowResponse) (response models.ShowResponse)
}

type Mapper struct{}

func (m Mapper) PhishNetResponseToShowResponse(pnResponse PNShowResponse) (response models.ShowResponse) {
	var songs []models.Song
	for _, data := range pnResponse.Data {
		songs = append(songs, models.Song{
			Title:     data.Song,
			TrackTime: data.Tracktime,
		})
	}

	response.Show = models.Show{
		Date:  pnResponse.Data[0].Showdate,
		Venue: pnResponse.Data[0].Venue,
		Songs: songs,
	}

	return response
}
