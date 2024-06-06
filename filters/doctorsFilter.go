package filters

import "gorm.io/gorm"

type DoctorsFilter struct {
	SpecialtyID string `in:"query=specialty"`
	Global      string `in:"query=global"`
}

func (doctorsFilter *DoctorsFilter) SpecialtyFilter(db *gorm.DB) *gorm.DB {
	return db.
		Where("specialty_id = ?", doctorsFilter.SpecialtyID)
}
