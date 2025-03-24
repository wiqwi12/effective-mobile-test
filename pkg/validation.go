package pkg

import "github.com/wiqwi12/effective-mobile-test/internal/domain/models/dto"

func IsEmpty(r dto.FilteredRequest) bool {
	return r.Title == "" && r.GroupName == "" && r.ReleaseDate == "" && r.Text == "" && r.Link == ""
}
