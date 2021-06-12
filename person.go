package churchtools

import (
	"encoding/json"
	"fmt"
)

type Email struct {
	Email          string
	IsDefault      bool
	ContactLabelID int
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

type Relationship struct {
	RelationshipTypeID   int
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
	FamilyStatusID             int
	WeddingDate                string
	StatusID                   int
	DateOfBaptism              string
	CanChat                    bool
	InvitationStatus           string
	ChatActive                 bool
	Meta                       MetaInfo
	IsArchived                 bool
	// TODO LatitudeLoose
	// TODO LongitudeLoose
	// TODO PrivacyPolicyAgreement
}

type PersonsResult struct {
	Data []Person
}

type RelationshipsResult struct {
	Data []Relationship
}

func (conn *Connector) GetPersons() ([]Person, error) {
	var person_list []Person

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
			person_list = append(person_list, person)
		}
	}

	return person_list, nil
}

func (conn *Connector) GetRelationships(id int) ([]Relationship, error) {
	url := fmt.Sprintf("persons/%d/relationships", id)
	content, err := conn.Get(url, true)
	if err != nil {
		return nil, err
	}

	var result RelationshipsResult
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
