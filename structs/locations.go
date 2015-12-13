package structs

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/owm-inc/weatherchecker-go/db"
)

type LocationEntry struct {
	Id                    bson.ObjectId `bson:"_id,omitempty" json:"objectid"`
	City_name             string        `json:"city_name"`
	Iso_country           string        `json:"iso_country"`
	Country_name          string        `json:"country_name"`
	Latitude              string        `json:"latitude"`
	Longitude             string        `json:"longitude"`
	Accuweather_id        string        `json:"accuweather_id"`
	Accuweather_city_name string        `json:"accuweather_city_name"`
	Gismeteo_id           string        `json:"gismeteo_id"`
	Gismeteo_city_name    string        `json:"gismeteo_city_name"`
}

func NewLocationEntry(city_name string, iso_country string, country_name string, latitude string, longitude string, accuweather_id string, accuweather_city_name string, gismeteo_id string, gismeteo_city_name string) LocationEntry {
	model := LocationEntry{City_name: city_name, Iso_country: iso_country, Country_name: country_name, Latitude: latitude, Longitude: longitude, Accuweather_id: accuweather_id, Accuweather_city_name: accuweather_city_name, Gismeteo_id: gismeteo_id, Gismeteo_city_name: gismeteo_city_name}
	model.Id = bson.NewObjectId()

	return model
}

type LocationTable struct {
	Database   *db.MongoDb
	Collection string
}

func (this *LocationTable) AddLocation(city_name string, iso_country string, country_name string, latitude string, longitude string, accuweather_id string, accuweather_city_name string, gismeteo_id string, gismeteo_city_name string) LocationEntry {
	newLocationEntry := NewLocationEntry(city_name, iso_country, country_name, latitude, longitude, accuweather_id, accuweather_city_name, gismeteo_id, gismeteo_city_name)
	this.Database.Insert(this.Collection, newLocationEntry)

	return newLocationEntry
}

func (this *LocationTable) RetrieveLocations() []LocationEntry {
	var result []LocationEntry
	this.Database.FindAll(this.Collection, &result)
	return result
}

func (this *LocationTable) RemoveLocation(location_id string) error {
	var err error
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
