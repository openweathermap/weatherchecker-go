package models

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/owm-inc/weatherchecker-go/db"
	"github.com/owm-inc/weatherchecker-go/util"
)

// LocationEntryBase represents the key fields of LocationEntry.
type LocationEntryBase struct {
	CityName    string `bson:"city_name" json:"city_name"`
	Slug        string `bson:"-" json:"slug"`
	IsoCountry  string `bson:"iso_country" json:"iso_country"`
	CountryName string `bson:"country_name" json:"country_name"`
	Latitude    string `bson:"latitude" json:"latitude"`
	Longitude   string `bson:"longitude" json:"longitude"`
}

// LocationEntry is the single location that will be queried for by Weather Checker.
type LocationEntry struct {
	DbEntryBase       `bson:",inline"`
	LocationEntryBase `bson:",inline"`
}

// NewLocationEntry makes a new location entry based on specified parameters.
func NewLocationEntry(
	cityName,
	isoCountry,
	countryName,
	latitude,
	longitude string) LocationEntry {
	model := LocationEntry{DbEntryBase{Id: bson.NewObjectId()}, LocationEntryBase{CityName: cityName, IsoCountry: isoCountry, CountryName: countryName, Latitude: latitude, Longitude: longitude}}

	return model
}

// LocationTable is a structure that acts as an interface between DB collection and Golang logic.
type LocationTable struct {
	Database   *db.MongoDb
	Collection string
}

// CreateLocation creates new location entry and inserts it into database.
func (c *LocationTable) CreateLocation(
	cityName,
	isoCountry,
	countryName,
	latitude,
	longitude string) (entry LocationEntry) {
	entry = NewLocationEntry(cityName, isoCountry, countryName, latitude, longitude)
	c.Database.Insert(c.Collection, entry)

	return entry
}

// ReadLocations returns all location entries in the database.
func (c *LocationTable) ReadLocations() []LocationEntry {
	var result []LocationEntry
	c.Database.FindAll(c.Collection, &result)

	output := make([]LocationEntry, len(result))
	for i, location := range result {
		location.Slug = util.MakeSlug(location.CityName)
		output[i] = location
	}
	return output
}

// UpdateLocation modifies location entry based on input parameters.
func (c *LocationTable) UpdateLocation(
	locationID,
	cityName,
	isoCountry,
	countryName,
	latitude,
	longitude string) (entry LocationEntry, err error) {
	b, idParseErr := db.GetObjectIDFromString(locationID)

	if idParseErr != nil {
		err = idParseErr
	} else {
		entry = NewLocationEntry(cityName, isoCountry, countryName, latitude, longitude)
		entry.Id = b
		err = c.Database.Update(c.Collection, b, entry)
	}
	return entry, err
}

// DeleteLocation removes location from the database.
func (c *LocationTable) DeleteLocation(locationID string) (err error) {
	b, idParseErr := db.GetObjectIDFromString(locationID)

	if idParseErr != nil {
		err = idParseErr
	} else {
		err = c.Database.Remove(c.Collection, b)
	}

	return err
}

// Clear removes all location entries from the database.
func (c *LocationTable) Clear() error {
	return c.Database.RemoveAll(c.Collection)
}

// NewLocationTable creates a new instance of LocationTable.
func NewLocationTable(dbInstance *db.MongoDb) LocationTable {
	var locations = LocationTable{Database: dbInstance, Collection: "Locations"}

	return locations
}
