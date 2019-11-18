package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

//struct Apply (payload request apply jobs)
type Apply struct {
	Nama        string
	JobId       int
	Age         int
	Gender      string
	PhoneNumber string
	Email       string
	Country     string
	Id          int
}

//struct Jobs, convert respons  https://www.themuse.com/api/public/jobs/:id,
type Jobs struct {
	Contents        string    `json:"contents"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	PublicationDate time.Time `json:"publication_date"`
	ShortName       string    `json:"short_name"`
	ModelType       string    `json:"model_type"`
	ID              int       `json:"id"`
	Locations       []struct {
		Name string `json:"name"`
	} `json:"locations"`
	Categories []struct {
		Name string `json:"name"`
	} `json:"categories"`
	Levels []struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
	} `json:"levels"`
	Tags []interface{} `json:"tags"`
	Refs struct {
		LandingPage string `json:"landing_page"`
	} `json:"refs"`
	Company struct {
		ID        int    `json:"id"`
		ShortName string `json:"short_name"`
		Name      string `json:"name"`
	} `json:"company"`
}

//convert respons https://restcountries.eu/rest/v2/name/{name}?fullText=true

type CountryStruct []struct {
	Name           string    `json:"name"`
	TopLevelDomain []string  `json:"topLevelDomain"`
	Alpha2Code     string    `json:"alpha2Code"`
	Alpha3Code     string    `json:"alpha3Code"`
	CallingCodes   []string  `json:"callingCodes"`
	Capital        string    `json:"capital"`
	AltSpellings   []string  `json:"altSpellings"`
	Region         string    `json:"region"`
	Subregion      string    `json:"subregion"`
	Population     int       `json:"population"`
	Latlng         []float64 `json:"latlng"`
	Demonym        string    `json:"demonym"`
	Area           float64   `json:"area"`
	Gini           float64   `json:"gini"`
	Timezones      []string  `json:"timezones"`
	Borders        []string  `json:"borders"`
	NativeName     string    `json:"nativeName"`
	NumericCode    string    `json:"numericCode"`
}

type ErrorCountry struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type refs struct {
	landingpage string
}

type company struct {
	id        string
	shortname string
	name      string
}

//declare var to save data apply
var applys []Apply

func ApplyController(c echo.Context) error {

	apply := Apply{}
	c.Bind(&apply)

	jobName := ValidateJobId(apply.JobId)

	//validate jobname, if job name null
	if jobName == "" {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Job Not Found",
		})
	}
	Countryname := ValidateCountry(apply.Country)
	// validate country
	if Countryname == "" {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Country Not Found",
		})
	}
	//ValidateName
	Nama := ValidateName(apply.Nama)
	if Nama == "" {
		return c.JSON(http.StatusNotAcceptable, map[string]interface{}{
			"message": "406 invalida data request",
		})
	}
	//ValidateEmaila
	Email := ValidateEmail(apply.Email)
	if Email == "" {
		return c.JSON(http.StatusNotAcceptable, map[string]interface{}{
			"message": "406 invalida data request",
		})
	}

	//increment UserId
	if len(applys) == 0 {
		apply.Id = 1
	} else {
		newId := applys[len(applys)-1].Id + 1
		apply.Id = newId
	}
	applys = append(applys, apply)

	//apply success return
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "succes apply",
		"Nama":    apply.Nama,
		"jobName": ValidateJobId(apply.JobId),
		"Age":     apply.Age,
		"email":   apply.Email,
		"Gender":  apply.Gender,
		"country": ValidateCountry(apply.Country),
		"id":      apply.Id,
	})
}

//validateJob
func ValidateJobId(JobId int) string {
	response, _ := http.Get("https://www.themuse.com/api/public/jobs/" + strconv.Itoa(JobId))

	responseData, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var Job Jobs

	//convet json to object / array
	json.Unmarshal([]byte(responseData), &Job)

	return Job.Name
}
func ValidateName(Nama string) string {
	return Nama
}

func ValidateEmail(Email string) string {
	return Email
}

//validate country
func ValidateCountry(Country string) string {
	//connect openapi
	responseCountry, _ := http.Get("https://www.restcountries.eu/rest/v2/name/" + Country + "?fullText=true")
	responseDataCountry, _ := ioutil.ReadAll(responseCountry.Body)
	defer responseCountry.Body.Close()

	var Countrys CountryStruct
	var Error ErrorCountry
	json.Unmarshal([]byte(responseDataCountry), &Countrys)
	json.Unmarshal([]byte(responseDataCountry), &Error)

	if Error.Message != "" {
		return ""
	}
	return Countrys[0].Name
}

//get user
func GetUsersController(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all users",
		"users":   applys,
	})
}

func GetUserController(c echo.Context) error {

	iddata, _ := strconv.Atoi(c.Param("id"))
	for _, value := range applys {
		if value.Id == iddata {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "success get data",
				"users":   value,
			})
		}
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"message": "data not found",
	})
}

func GetUserCountry(c echo.Context) error {
	var locationjobs string
	locations := strings.ToLower(c.QueryParam("locations"))
	var Kerjaan Jobs1
	response, _ := http.Get("https://www.themuse.com/api/public/jobs?page=1")

	responseData, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal(responseData, &Kerjaan)
	// fmt.Println(Kerjaan)
	for _, value := range Kerjaan.Results {
		if (len(value.Locations)) < 1 {
			continue
		} else {
			locationjobs = strings.ToLower(value.Locations[0].Name)
		}
		if strings.Contains(locationjobs, locations) {

			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "success get data",
				"users":   value,
			})
		}
	}
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"message": "data not found",
	})
}

// fmt.Println("======================== LIST COUNTRY =================================")

type AutoGenerated []struct {
	Name           string   `json:"name"`
	TopLevelDomain []string `json:"topLevelDomain"`
}

func listCountry(c echo.Context) error {
	response, _ := http.Get("https://restcountries.eu/rest/v2/all")

	responseData, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var data AutoGenerated
	json.Unmarshal(responseData, &data)
	// for _,allCountry := range data{
	// 	fmt.Println("Name   : ", allCountry.Name)
	// 	fmt.Println("TopLevelDomain : ", allCountry.TopLevelDomain)
	// }

	return c.JSON(http.StatusOK, data)
}
