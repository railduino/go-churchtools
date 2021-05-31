package churchtools

import (
	"encoding/json"
	"fmt"
)

type Email struct {
	Email          string
	IsDefault      bool
	ContactLabelId int
}

type DomainAttributes struct {
	FirstName string
	LastName  string
	GUID      string
}

type Relative struct {
	Title            string
	DomainType       string
	DomainIdentifier string
	ApiUrl           string
	FrontendUrl      string
	ImageUrl         string
	DomainAttributes DomainAttributes
}

type Relation struct {
	RelationshipTypeId   int
	RelationshipName     string
	DegreeOfRelationship string
	Relative             Relative
}

type Person struct {
	ID                         int
	GUID                       string
	SecurityLevelForPerson     int
	EditSecurityLevelForPerson int
	FirstName                  string
	LastName                   string
	Nickname                   string
	Street                     string
	AddressAddition            string
	Zip                        string
	City                       string
	Country                    string
	Latitude                   float64
	Longitude                  float64
	PhonePrivate               string
	PhoneWork                  string
	Mobile                     string
	BirthName                  string
	Birthday                   string
	ImageUrl                   string
	FamilyImageUrl             string
	Email                      string
	Emails                     []Email
	familyStatusId             int
	weddingDate                string
	statusId                   int
	dateOfBaptism              string
	canChat                    bool
	invitationStatus           string
	chatActive                 bool
	Meta                       MetaInfo
	isArchived                 bool
	Relations                  []Relation `json:"-"`
	// TODO LatitudeLoose
	// TODO LongitudeLoose
	// TODO privacyPolicyAgreement
}

type PersonsResult struct {
	Data []Person
}

var (
	Persons       []Person
	PersonIdMap   = make(map[int]*Person)
	PersonGuidMap = make(map[string]*Person)
)

func (conn *Connector) GetPersons() ([]Person, error) {
	for page := 1; ; page++ {
		url := fmt.Sprintf("persons?page=%d&limit=100", page)
		content, err := conn.Get(url, true)
		if err != nil {
			return nil, err
		}

		var result PersonsResult
		if err := json.Unmarshal(content, &result); err != nil {
			return nil, err
		}

		if len(result.Data) == 0 {
			break
		}

		for _, person := range result.Data {
			Persons = append(Persons, person)
			PersonIdMap[person.ID] = &person
			PersonGuidMap[person.GUID] = &person
		}
	}

	return Persons, nil
}
