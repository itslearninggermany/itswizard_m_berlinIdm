package itswizard_m_berlinIdm

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

type LusdUser struct {
	PersonUID            string `json:"personUID"`
	BenutzerVorname      string `json:"benutzerVorname"`
	BenutzerRufname      string `json:"benutzerRufname"`
	BenutzerNachname     string `json:"benutzerNachname"`
	BenutzerTitel        string `json:"benutzerTitel"`
	BenutzerGeburtsdatum string `json:"benutzerGeburtsdatum"`
	BenutzerGlobalRolle  string `json:"benutzerGlobalRolle"`
	BenutzerAktionCode   string `json:"benutzerAktionCode"`
	BenutzerDtAktion     string `json:"benutzerDtAktion"`
	BenutzerSchuleListe  []struct {
		SchuleUID           string `json:"schuleUID"`
		BenutzerStatus      string `json:"benutzerStatus"`
		BenutzerSchuleRolle string `json:"benutzerSchuleRolle"`
		BenutzerKlasseListe []struct {
			KlasseUID           string `json:"klasseUID"`
			BenutzerKlasseRolle string `json:"benutzerKlasseRolle"`
		} `json:"benutzerKlasseListe"`
		BenutzerKursListe []struct {
			KursUID           string `json:"kursUID"`
			BenutzerKursRolle string `json:"benutzerKursRolle"`
		} `json:"benutzerKursListe"`
	} `json:"benutzerSchuleListe"`
	BenutzerBezugspersonListe []struct {
		BezugspersonPersonUID string `json:"bezugspersonPersonUID"`
		BezugspersonTyp       string `json:"bezugspersonTyp"`
	} `json:"benutzerBezugspersonListe"`
}

func personsParse(out []byte, err error) (lusdUsers []LusdUser, err2 error) {
	if err != nil {
		return nil, err
	}
	err2 = json.Unmarshal(out, &lusdUsers)
	return lusdUsers, err2
}

/*
func (p *blusdConnection) GetAllPersons() (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(time.Now(), "", users, true))
}
*/

// Wenn die Zahl negativ ist, wird der 01.01.2022 00:00:01.000 genutzt
func (p *BlusdConnection) GetAllPersons(anzahlAnVergangenenMinuten int) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPIPerson(anzahlAnVergangenenMinuten, "", users, false))
}

/*
Returns all Persons from a school
*/
func (p *BlusdConnection) GetAllPersonsFromSchool(sid string) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(time.Now(), sid, users, false))
}

/*
Returns all new Persons from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllPersonsFromSchoolScince(sid string, minutes time.Duration) (lusdUsers []LusdUser, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return personsParse(p.callAPI(newT, sid, users, true))
}

/*
Returns all Persons from a school created since a given time
*/
func (p *BlusdConnection) GetAllPersonsFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(t, sid, users, true))
}

/*
Returns all Persons created since a given time
*/
func (p *BlusdConnection) GetAllPersonsScinceAGivenTime(t time.Time) (lusdUsers []LusdUser, err error) {
	return personsParse(p.callAPI(t, "", users, true))
}

/*
Get details from one given Person
*/
func (p *BlusdConnection) GetDetailOfPerson(personID string) (lusdUsers LusdUser, err error) {
	tmp, err := personsParse(p.callAPI(time.Now(), "", user+personID, false))
	if err != nil {
		if len(tmp) < 1 {
			erro := err.Error() + "The person with the personID " + personID + " does not exist!"
			return lusdUsers, errors.New(erro)
		}
		return lusdUsers, err
	}
	if len(tmp) < 1 {
		return lusdUsers, errors.New("The person with the personID " + personID + " does not exist!")
	}
	return tmp[0], err
}

func (p *LusdUser) GetAllSchools() (out string, err error) {
	m := []string{}
	for _, k := range p.BenutzerSchuleListe {
		m = append(m, k.SchuleUID)
	}
	b, err := json.Marshal(m)
	return string(b), err
	return out, err
}

func (p *LusdUser) GetChildParentRelation() (out string, err error) {
	m := make(map[string]string)
	for _, k := range p.BenutzerBezugspersonListe {
		m[k.BezugspersonPersonUID] = k.BezugspersonTyp
	}
	b, err := json.Marshal(m)
	return string(b), err
	return out, err
}

func (p *LusdUser) GetAllClasses() (out string, err error) {
	m := make(map[string]string)

	for _, k := range p.BenutzerSchuleListe {
		for _, k2 := range k.BenutzerKlasseListe {
			m[k2.KlasseUID] = k2.BenutzerKlasseRolle
		}
	}
	b, err := json.Marshal(m)
	return string(b), err
}

func (p *LusdUser) GetAllCourses() (out string, err error) {
	m := make(map[string]string)

	for _, k := range p.BenutzerSchuleListe {
		for _, k2 := range k.BenutzerKursListe {
			m[k2.KursUID] = k2.BenutzerKursRolle
		}
	}
	b, err := json.Marshal(m)
	return string(b), err
}
