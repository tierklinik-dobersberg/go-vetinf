package vetinf

import "fmt"

type Meta struct {
	Deleted bool
}

// Customer represents customer data stored in the vetkldat.dbf file
type Customer struct {
	Meta                Meta   `dbf:"-"`
	ID                  int    `dbf:"knr" json:"id,omitempty" bson:"id,omitempty"`
	Group               string `dbf:"gruppe" json:"group,omitempty" bson:"group,omitempty"`
	Name                string `dbf:"name" json:"name,omitempty" bson:"name:omitempty"`
	Firstname           string `dbf:"vorname" json:"firstname,omitempty" bson:"firstname,omitempty"`
	Titel               string `dbf:"titel" json:"title,omitempty" bson:"title,omitempty"`
	Street              string `dbf:"strasse" json:"street,omitempty" bson:"street,omitempty"`
	CityCode            int    `dbf:"plz" json:"cityCode,omitempty" bson:"cityCode,omitempty"`
	City                string `dbf:"ort" json:"city,omitempty" bson:"city,omitempty"`
	Phone               string `dbf:"telefon" json:"phone,omitempty" bson:"phone,omitempty"`
	Extra               string `dbf:"zusatz" json:"extra,omitempty" bson:"extra,omitempty"`
	Salutation          string `dbf:"anrede" json:"salutation,omitempty" bson:"salutation,omitempty"`
	MobilePhone1        string `dbf:"handyx1" json:"mobilePhone1,omitempty" bson:"mobilePhone1,omitempty"`
	MobilePhone2        string `dbf:"handyx2" json:"mobilePhone2,omitempty" bson:"mobilePhone2,omitempty"`
	Phone2              string `dbf:"telefonx1" json:"phone2,omitempty" bson:"phone2,omitempty"`
	Mail                string `dbf:"emailx1" json:"mail,omitempty" bson:"mail,omitempty"`
	SecondaryID         string `dbf:"knr2" json:"secondaryId,omitempty" bson:"secondaryId,omitempty"`
	VaccinationReminder bool   `dbf:"impfung" json:"vaccinationReminder,omitempty" bson:"vaccinationReminder,omitempty"`
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer{id:%d name:%q}", c.ID, c.Name)
}

type SmallAnimalRecord struct {
	Meta Meta `dbf:"-"`

	CustomerID    int    `dbf:"knr" json:"knr" bson:"knr"`
	Size          string `dbf:"grossklein" json:"size" bson:"size"`
	Species       string `dbf:"tierart" json:"species" bson:"species"`
	Breed         string `dbf:"rasse" json:"breed" bson:"breed"`
	Gender        string `dbf:"geschlecht" json:"gender" bson:"gender"`
	Name          string `dbf:"name" json:"name" bson:"name"`
	Birthday      string `dbf:"gebdat" json:"birthday" bson:"birthday"`
	SpecialDetail string `dbf:"besonderes" json:"specialDetail" bson:"specialDetail"`
	AnimalID      string `dbf:"tilfdnr" json:"animalId" bson:"animalId"`
	Extra1        string `dbf:"mehr1" json:"extra1" bson:"extra1"`
	Extra2        string `dbf:"mehr2" json:"extra2" bson:"extra2"`
	Extra3        string `dbf:"mehr3" json:"extra3" bson:"extra3"`
	Extra4        string `dbf:"mehr4" json:"extra4" bson:"extra4"`
	Extra5        string `dbf:"mehr5" json:"extra5" bson:"extra5"`
	Extra6        string `dbf:"mehr6" json:"extra6" bson:"extra6"`
	Extra7        string `dbf:"mehr7" json:"extra7" bson:"extra7"`
	Extra8        string `dbf:"mehr8" json:"extra8" bson:"extra8"`
	Extra9        string `dbf:"mehr9" json:"extra9" bson:"extra9"`
	Extra10       string `dbf:"mehr10" json:"extra10" bson:"extra10"`
	Color         string `dbf:"farbe" json:"color" bson:"color"`
	ChipNumber    string `dbf:"chipnr" json:"chipNumber" bson:"chipNumber"`
}
