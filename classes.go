package itswizard_m_berlinIdm

import (
	"encoding/json"
	"log"
	"time"
)

type LusdClass struct {
	KlasseDtGeaendertAm string `json:"klasseDtGeaendertAm"`
	KlasseName          string `json:"klasseName"`
	KlasseUID           string `json:"klasseUID"`
	SchuleUID           string `json:"schuleUID"`
}

func classParse(out []byte, err error) (lusdClasses []LusdClass, err2 error) {
	if err != nil {
		return lusdClasses, err
	}
	err2 = json.Unmarshal(out, &lusdClasses)
	return lusdClasses, err2
}

/*
Returns all Schools from lusd
*/
func (p *BlusdConnection) GetAllClasses() (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(time.Now(), "", class, false))
}

/*
Returns all new Courses created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllClassessCreatedScinse(minutes time.Duration) (lusdClasses []LusdClass, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return classParse(p.callAPI(newT, "", class, true))
}

/*
Returns all Courses from a school
*/
func (p *BlusdConnection) GetAllClassesFromSchool(sid string) (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(time.Now(), sid, class, false))
}

/*
Returns all new Courses from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllClassesFromSchoolScince(sid string, minutes time.Duration) (lusdClasses []LusdClass, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return classParse(p.callAPI(newT, sid, class, true))
}

/*
Returns all Courses from a school created since a given time
*/
func (p *BlusdConnection) GetAllClassesFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(t, sid, class, true))
}

/*
Returns all Courses created since a given time
*/
func (p *BlusdConnection) GetAllClassesScinceAGivenTime(t time.Time) (lusdClasses []LusdClass, err error) {
	return classParse(p.callAPI(t, "", class, true))
}
