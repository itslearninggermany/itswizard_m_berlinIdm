package itswizard_m_berlinIdm

import (
	"encoding/json"
	"log"
	"time"
)

type LusdCourse struct {
	KursBezeichnung   string `json:"kursBezeichnung"`
	KursDtGeaendertAm string `json:"kursDtGeaendertAm"`
	KursFach          string `json:"kursFach"`
	KursJahrgang      string `json:"kursJahrgang"`
	KursSchulform     string `json:"kursSchulform"`
	KursUID           string `json:"kursUID"`
	SchuleUID         string `json:"schuleUID"`
}

func courseParse(out []byte, err error) (lusdCourses []LusdCourse, err2 error) {
	if err != nil {
		return nil, err
	}
	err2 = json.Unmarshal(out, &lusdCourses)
	return lusdCourses, err2
}

/*
Returns all Courses from lusd
*/
func (p *BlusdConnection) GetAllCourses() (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(time.Now(), "", course, false))
}

/*
Returns all new Courses created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllCoursesCreatedScinse(minutes time.Duration) (lusdCourses []LusdCourse, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return courseParse(p.callAPI(newT, "", course, true))
}

/*
Returns all Courses from a school
*/
func (p *BlusdConnection) GetAllCoursesFromSchool(sid string) (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(time.Now(), sid, course, false))
}

/*
Returns all new Courses from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllCoursesFromSchoolScince(sid string, minutes time.Duration) (lusdCourses []LusdCourse, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return courseParse(p.callAPI(newT, sid, course, true))
}

/*
Returns all Courses from a school created since a given time
*/
func (p *BlusdConnection) GetAllCoursesFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(t, sid, course, true))
}

/*
Returns all Courses created since a given time
*/
func (p *BlusdConnection) GetAllCoursesScinceAGivenTime(t time.Time) (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(t, "", course, true))
}
