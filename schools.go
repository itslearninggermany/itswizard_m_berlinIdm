package itswizard_m_berlinIdm

import (
	"encoding/json"
	"log"
	"time"
)

type LusdSchool struct {
	SchuleBezirkName     string   `json:"schuleBezirkName"`
	SchuleBezirkNummer   string   `json:"schuleBezirkNummer"`
	SchuleDtGeaendertAm  string   `json:"schuleDtGeaendertAm"`
	SchuleName           string   `json:"schuleName"`
	SchuleNummer         string   `json:"schuleNummer"`
	SchuleSchulformListe []string `json:"schuleSchulformListe"`
	SchuleTyp            string   `json:"schuleTyp"`
	SchuleUID            string   `json:"schuleUID"`
}

func schoolParse(out []byte, err error) (lusdSchools []LusdSchool, err2 error) {
	if err != nil {
		return lusdSchools, err
	}
	err2 = json.Unmarshal(out, &lusdSchools)
	return lusdSchools, err2
}

/*
Returns all Schools from lusd
*/
func (p *BlusdConnection) GetAllSchools() (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(time.Now(), "", school, false))
}

/*
Returns all new Courses created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllSchoolsCreatedScinse(minutes time.Duration) (lusdSchools []LusdSchool, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return schoolParse(p.callAPI(newT, "", school, true))
}

/*
Returns all Courses from a school
*/
func (p *BlusdConnection) GetAllSchoolsFromSchool(sid string) (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(time.Now(), sid, school, false))
}

/*
Returns all new Courses from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllSchoolsFromSchoolScince(sid string, minutes time.Duration) (lusdSchools []LusdSchool, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return schoolParse(p.callAPI(newT, sid, school, true))
}

/*
Returns all Courses from a school created since a given time
*/
func (p *BlusdConnection) GetAllSchoolsFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(t, sid, school, true))
}

/*
Returns all Courses created since a given time
*/
func (p *BlusdConnection) GetAllSchoolsScinceAGivenTime(t time.Time) (lusdSchools []LusdSchool, err error) {
	return schoolParse(p.callAPI(t, "", school, true))
}
