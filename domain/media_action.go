package domain

func (m MediaMongo) ToMedia() Media {
	return Media(m)
}

func (m Media) ToMediaMongo() MediaMongo {
	return MediaMongo(m)
}
