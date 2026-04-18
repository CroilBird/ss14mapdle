package models

type MapRedirect struct {
	Id    uint64
	MapID uint64
	Map   Map
}

func init() {
	Models = append(Models, MapRedirect{})
}
