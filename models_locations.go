package main

import (
	"strconv"

	"github.com/skybon/mgoHelpers"
	"gopkg.in/mgo.v2/bson"
)

// LocationEntryBase represents the key fields of LocationEntry.
type LocationEntryBase struct {
	CityName    string `bson:"city_name" json:"city_name"`
	IsoCountry  string `bson:"iso_country" json:"iso_country"`
	CountryName string `bson:"country_name" json:"country_name"`
	Latitude    string `bson:"latitude" json:"latitude"`
	Longitude   string `bson:"longitude" json:"longitude"`
}

// LocationEntry is the single location that will be queried for by Weather Checker.
type LocationEntry struct {
	mgoHelpers.DbEntryBase `bson:",inline"`
	LocationEntryBase      `bson:",inline"`
	Slug                   string `bson:"slug" json:"slug"`
}

// NewLocationEntry makes a new location entry based on specified parameters.
func NewLocationEntry(entryBase LocationEntryBase) LocationEntry {
	var entry = LocationEntry{LocationEntryBase: entryBase}

	entry.SetBsonID(bson.NewObjectId())

	return entry
}

func makeUniqueSlug(c *mgoHelpers.MongoCollection, entry LocationEntry) string {
	var slug string
	for i := 0; ; i++ {
		var newSlug string
		if i == 0 {
			newSlug = MakeSlug(entry.CityName)
		} else {
			newSlug = MakeSlug(entry.CityName + "_" + strconv.FormatInt(int64(i), 64))
		}

		var existingSlugs []LocationEntry
		c.Database.Find(c.Collection, map[string]interface{}{"slug": newSlug, "_id": map[string]interface{}{"$ne": entry.BsonID()}}, &existingSlugs)

		if len(existingSlugs) == 0 {
			slug = newSlug
			break
		}
	}

	return slug
}

type LocationCollection struct {
	base *mgoHelpers.MongoCollection
}

func (c *LocationCollection) Create(entryBase LocationEntryBase) (entry LocationEntry, err error) {
	abstractEntry, dbErr := c.base.Create(entryBase)

	if dbErr != nil {
		return entry, dbErr
	}

	locEntryP, assertOk := abstractEntry.(*LocationEntry)
	if !assertOk {
		return entry, MalformedEntry
	}

	entry = *locEntryP
	return entry, nil
}

func (c *LocationCollection) Read(entryID string) (entry LocationEntry, exists bool) {
	exists = c.base.Read(entryID, &entry)

	return entry, exists
}

func (c *LocationCollection) ReadAll() (output []LocationEntry, err error) {
	err = c.base.ReadAll(&output)

	return output, err
}

func (c *LocationCollection) Update(locationID string, entryBase LocationEntryBase) (LocationEntry, error) {
	var err error

	abstractEntry, updErr := c.base.Update(locationID, entryBase)

	if updErr != nil {
		return LocationEntry{}, updErr
	}
	entry, matches := abstractEntry.(*LocationEntry)

	if !matches {
		err = MalformedEntry
	}

	return *entry, err

}

func (c *LocationCollection) Delete(locationID string) error { return c.base.Delete(locationID) }

func (c *LocationCollection) DeleteAll() error { return c.base.DeleteAll() }

func NewLocationCollection(dbInstance *mgoHelpers.MongoDb) *LocationCollection {
	var coll LocationCollection
	coll.base = mgoHelpers.NewMongoCollection(dbInstance, "Locations")
	coll.base.SetFactoryFunc(func(coll *mgoHelpers.MongoCollection, factoryArgs interface{}) mgoHelpers.MongoEntry {
		entryBase := factoryArgs.(LocationEntryBase)

		entry := NewLocationEntry(entryBase)
		entry.Slug = makeUniqueSlug(coll, entry)

		return &entry
	})
	return &coll
}
