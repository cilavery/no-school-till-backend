package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cilavery/no-school-till-backend/cmd/main.go/utils"
	"github.com/gorilla/mux"
)

var client = &http.Client{}
var baseURL = "https://developers.teachable.com/v1"

type Users struct {
	Users []user `json:"users"`
	Meta  meta   `json:"meta"`
}

type Courses struct {
	Courses []struct {
		ID          int    `json:"id"`
		Description any    `json:"description"`
		Name        string `json:"name"`
		Heading     string `json:"heading"`
		IsPublished bool   `json:"is_published"`
		ImageURL    string `json:"image_url"`
	} `json:"courses"`
	Meta meta `json:"meta"`
}

type Enrollments struct {
	Enrollments []struct {
		UserID          int       `json:"user_id"`
		EnrolledAt      time.Time `json:"enrolled_at"`
		CompletedAt     any       `json:"completed_at"`
		PercentComplete int       `json:"percent_complete"`
		ExpiresAt       any       `json:"expires_at"`
	} `json:"enrollments"`
	Meta meta `json:"meta"`
}

type meta struct {
	Total         int `json:"total"`
	Page          int `json:"page"`
	From          int `json:"from"`
	To            int `json:"to"`
	PerPage       int `json:"per_page"`
	NumberOfPages int `json:"number_of_pages"`
}

type user struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Controller struct{}

// With time would refactor to use an interface so Courses and Users would both have the same interface

func NewController() *Controller {
	return &Controller{}
}

var allStudents = make(map[int]user)

func (c *Controller) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s/users", baseURL)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}

		setHeaders(req)

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer res.Body.Close()

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		var users Users

		if err = json.Unmarshal(responseBody, &users); err != nil {
			// TODO: handle by response status codes and send
			// proper error statuses
			utils.SendError(w, http.StatusBadRequest, utils.Error{Message: "response error"})
			return
		}

		for _, u := range users.Users {
			// check if user_id exists in the student map, if not then add it
			if _, ok := allStudents[u.ID]; !ok {
				var student user
				student.Name = u.Name
				student.Email = u.Email
				allStudents[u.ID] = student
			}
		}
	}
}

func (c *Controller) GetAllCourses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("%s/courses", baseURL)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}

		setHeaders(req)

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer res.Body.Close()

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		var courses Courses

		if err = json.Unmarshal(responseBody, &courses); err != nil {
			// TODO: handle by response status codes and send
			// proper error statuses
			utils.SendError(w, http.StatusBadRequest, utils.Error{Message: "response error"})
			return
		}

		utils.SendSuccess(w, courses)
	}
}

func (c *Controller) GetAllEnrollments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["course_id"]
		url := fmt.Sprintf("%s/courses/%s/enrollments", baseURL, id)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
		}

		setHeaders(req)

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer res.Body.Close()

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		var enrollments Enrollments

		if err = json.Unmarshal(responseBody, &enrollments); err != nil {
			// TODO: handle by response status codes and send
			// proper error statuses
			utils.SendError(w, http.StatusBadRequest, utils.Error{Message: "response error"})
			return
		}

		enrolledStudents := []user{}
		// doesn't handle if a user_id doesn't exist in the students map
		for _, enrolled := range enrollments.Enrollments {
			studentData := allStudents[enrolled.UserID]
			studentData.ID = enrolled.UserID
			enrolledStudents = append(enrolledStudents, studentData)
		}

		utils.SendSuccess(w, enrolledStudents)
	}
}

func setHeaders(r *http.Request) *http.Request {
	apiKey := os.Getenv("API_KEY")
	r.Header.Add("apiKey", apiKey)
	return r
}
