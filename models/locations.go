package models

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/owm-inc/weatherchecker-go/db"
)

type LocationEntryBase struct {
	City_name    string `json:"city_name"`
	Iso_country  string `json:"iso_country"`
	Country_name string `json:"country_name"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
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
	longitude string) LocationEntry {
	model := LocationEntry{DbEntryBase{Id: bson.NewObjectId()}, LocationEntryBase{City_name: city_name, Iso_country: iso_country, Country_name: country_name, Latitude: latitude, Longitude: longitude}}

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
	longitude string) (entry LocationEntry) {
	entry = NewLocationEntry(city_name, iso_country, country_name, latitude, longitude)
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
	longitude string) (entry LocationEntry, err error) {
	b, idParseErr := db.GetObjectIDFromString(location_id)

	if idParseErr != nil {
		err = idParseErr
	} else {
		entry = NewLocationEntry(city_name, iso_country, country_name, latitude, longitude)
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
