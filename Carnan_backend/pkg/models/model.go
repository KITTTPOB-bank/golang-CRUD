package models

import "mime/multipart"

type Cardata struct {
	Name      string `bson:"name"`
	Brand     string `bson:"brand"`
	Year      int    `bson:"year"`
	Desc      string `bson:"desc"`
	Car_image string `bson:"car_image"`
}
type Carname struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	Year  int    `json:"year"`
}
type Carpart struct {
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Year      string `json:"year"`
	Type      int    `json:"type"`
	IDcarpart []int  `json:"selected_ids"`
}
type Dropcarpart struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	Year  string `json:"year"`
	Type  int    `json:"type"`
	Delid int    `json:"delid"`
}
type Formdatacar struct {
	Name      string                `form:"name"`
	Brand     string                `form:"brand"`
	Year      int                   `form:"year"`
	Desc      string                `form:"desc"`
	Car_image *multipart.FileHeader `form:"car_image"`
}
type Carpartdataall struct {
	Frontbumper_list []int `bson:"frontbumper_list"`
	Rearbumper_list  []int `bson:"rearbumper_list"`
	Grille_list      []int `bson:"grille_list"`
	Headlamp_list    []int `bson:"headlamp_list"`
	Backuplamp_list  []int `bson:"backuplamp_list"`
	Mirror_list      []int `bson:"mirror_list"`
	Door_list        []int `bson:"door_list"`
}
type Addcarpart struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Price string `json:"price"`
	Type  int    `json:"type"`
}

type FormJsonfilecarpart struct {
	FileJson *multipart.FileHeader `form:"file"`
}
type FormJsonfileMycarpart struct {
	FileJson  *multipart.FileHeader `form:"file"`
	Name      string                `form:"name"`
	Brand     string                `form:"brand"`
	Year      int                   `form:"year"`
	Typecheck int                   `form:"typecheck"`
}
type Cariamge struct {
	Name      string `form:"name"`
	Brand     string `form:"brand"`
	Year      int    `form:"year"`
	Typecheck int    `form:"typecheck"`
}
