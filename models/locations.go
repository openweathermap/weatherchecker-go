package models

import (
	"strconv"

	"github.com/skybon/semaphore"
	"gopkg.in/mgo.v2/bson"

	"github.com/owm-inc/weatherchecker-go/db"
	"github.com/owm-inc/weatherchecker-go/util"
)

// LocationEntryBase represents the key fields of LocationEntry.
type LocationEntryBase struct {
	CityName    string `bson:"city_name" json:"city_name"`
	Slug        string `bson:"slug" json:"slug"`
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
func NewLocationEntry(cityName, isoCountry, countryName, latitude, longitude string) LocationEntry {
	model := LocationEntry{DbEntryBase{Id: bson.NewObjectId()}, LocationEntryBase{CityName: cityName, IsoCountry: isoCountry, CountryName: countryName, Latitude: latitude, Longitude: longitude}}

	return model
}

// LocationTable is a structure that acts as an interface between DB collection and Golang logic.
type LocationTable struct {
	Database   *db.MongoDb
	Collection string
	semaphore  semaphore.Semaphore
}

func (c *LocationTable) makeUniqueSlug(entry LocationEntry) string {
	var slug string
	for i := 0; ; i++ {
		var newSlug string
		if i == 0 {
			newSlug = util.MakeSlug(entry.CityName)
		} else {
			newSlug = util.MakeSlug(entry.CityName + "_" + strconv.FormatInt(int64(i), 64))
		}

		var existingSlugs []LocationEntry
		c.Database.Find(c.Collection, map[string]interface{}{"slug": newSlug, "_id": map[string]interface{}{"$ne": entry.Id}}, &existingSlugs)

		if len(existingSlugs) == 0 {
			slug = newSlug
			break
		}
	}

	return slug
}

func (c *LocationTable) createLocationCore(cityName, isoCountry, countryName, latitude, longitude string) LocationEntry {
	entry := NewLocationEntry(cityName, isoCountry, countryName, latitude, longitude)
	slug := c.makeUniqueSlug(entry)
	entry.Slug = slug
	c.Database.Insert(c.Collection, entry)

	return entry
}

// CreateLocation creates new location entry and inserts it into database.
func (c *LocationTable) CreateLocation(cityName, isoCountry, countryName, latitude, longitude string) (entry LocationEntry) {
	c.semaphore.Exec(func() { entry = c.createLocationCore(cityName, isoCountry, countryName, latitude, longitude) })

	return entry
}

// ReadLocations returns all location entries in the database.
func (c *LocationTable) ReadLocations() []LocationEntry {
	var result []LocationEntry
	c.Database.FindAll(c.Collection, &result)

	output := make([]LocationEntry, len(result))
	for i, location := range result {
		output[i] = location
	}
	return output
}

// UpdateLocation modifies location entry based on input parameters.
func (c *LocationTable) UpdateLocation(locationID, cityName, isoCountry, countryName, latitude, longitude string) (entry LocationEntry, status error) {
	c.semaphore.Exec(func() {
		b, idParseErr := db.GetObjectIDFromString(locationID)

		var err error
		var newEntry LocationEntry

		if idParseErr != nil {
			err = idParseErr
		} else {
			newEntry = NewLocationEntry(cityName, isoCountry, countryName, latitude, longitude)
			newEntry.Id = b
			newEntry.Slug = c.makeUniqueSlug(newEntry)
			err = c.Database.Update(c.Collection, b, newEntry)
		}
		entry = newEntry
		status = err
	})
	return entry, status
}

// DeleteLocation removes location from the database.
func (c *LocationTable) DeleteLocation(locationID string) (status error) {
	c.semaphore.Exec(func() {
		var err error
		b, idParseErr := db.GetObjectIDFromString(locationID)

		if idParseErr != nil {
			err = idParseErr
		} else {
			err = c.Database.Remove(c.Collection, b)
		}
		status = err
	})

	return status
}

// Clear removes all location entries from the database.
func (c *LocationTable) Clear() (status error) {
	c.semaphore.Exec(func() { status = c.Database.RemoveAll(c.Collection) })

	return status
}

// NewLocationTable creates a new instance of LocationTable.
func NewLocationTable(dbInstance *db.MongoDb) LocationTable {
	return LocationTable{Database: dbInstance, Collection: "Locations", semaphore: semaphore.MakeSemaphore(1)}
}
