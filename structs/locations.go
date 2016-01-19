package structs

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/owm-inc/weatherchecker-go/db"
)

type LocationEntryBase struct {
	City_name             string `json:"city_name"`
	Iso_country           string `json:"iso_country"`
	Country_name          string `json:"country_name"`
	Latitude              string `json:"latitude"`
	Longitude             string `json:"longitude"`
	Accuweather_id        string `json:"accuweather_id"`
	Accuweather_city_name string `json:"accuweather_city_name"`
	Gismeteo_id           string `json:"gismeteo_id"`
	Gismeteo_city_name    string `json:"gismeteo_city_name"`
	Yandex_id             string `json:"yandex_id"`
}

type LocationEntry struct {
	DbEntryBase       `bson:",inline"`
	LocationEntryBase `bson:",inline"`
}

func NewLocationEntry(
	city_name,
	iso_country,
	country_name,
	latitude,
	longitude,
	accuweather_id,
	accuweather_city_name,
	gismeteo_id,
	gismeteo_city_name,
	yandex_id string) LocationEntry {
	model := LocationEntry{DbEntryBase{Id: bson.NewObjectId()}, LocationEntryBase{City_name: city_name, Iso_country: iso_country, Country_name: country_name, Latitude: latitude, Longitude: longitude, Accuweather_id: accuweather_id, Accuweather_city_name: accuweather_city_name, Gismeteo_id: gismeteo_id, Gismeteo_city_name: gismeteo_city_name, Yandex_id: yandex_id}}

	return model
}

type LocationTable struct {
	Database   *db.MongoDb
	Collection string
}

func (this *LocationTable) CreateLocation(
	city_name,
	iso_country,
	country_name,
	latitude,
	longitude,
	accuweather_id,
	accuweather_city_name,
	gismeteo_id,
	gismeteo_city_name,
	yandex_id string) (entry LocationEntry) {
	entry = NewLocationEntry(city_name, iso_country, country_name, latitude, longitude, accuweather_id, accuweather_city_name, gismeteo_id, gismeteo_city_name, yandex_id)
	this.Database.Insert(this.Collection, entry)

	return entry
}

func (this *LocationTable) ReadLocations() (result []LocationEntry) {
	this.Database.FindAll(this.Collection, &result)
	return result
}

func (this *LocationTable) UpdateLocation(
	location_id,
	city_name,
	iso_country,
	country_name,
	latitude,
	longitude,
	accuweather_id,
	accuweather_city_name,
	gismeteo_id,
	gismeteo_city_name,
	yandex_id string) (entry LocationEntry, err error) {
	b, idParseErr := db.GetObjectIDFromString(location_id)

	if idParseErr != nil {
		err = idParseErr
	} else {
		entry = NewLocationEntry(city_name, iso_country, country_name, latitude, longitude, accuweather_id, accuweather_city_name, gismeteo_id, gismeteo_city_name, yandex_id)
		entry.Id = b
		err = this.Database.Update(this.Collection, b, entry)
	}
	return entry, err
}

func (this *LocationTable) DeleteLocation(location_id string) (err error) {
	b, idParseErr := db.GetObjectIDFromString(location_id)

	if idParseErr != nil {
		err = idParseErr
	} else {
		err = this.Database.Remove(this.Collection, b)
	}

	return err
}

func (this *LocationTable) Clear() error {
	return this.Database.RemoveAll(this.Collection)
}

func NewLocationTable(db_instance *db.MongoDb) LocationTable {
	var locations = LocationTable{Database: db_instance, Collection: "Locations"}

	return locations
}
