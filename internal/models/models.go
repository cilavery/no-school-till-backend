package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/cilavery/no-school-till-backend/cmd/main.go/utils"
)

type Users struct {
	Users []User `json:"users"`
	Meta  Meta   `json:"meta"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Courses struct {
	Courses []Course `json:"courses"`
	Meta    Meta     `json:"meta"`
}

type Course struct {
	ID          int    `json:"id"`
	Description any    `json:"description"`
	Name        string `json:"name"`
	Heading     string `json:"heading"`
	IsPublished bool   `json:"is_published"`
	ImageURL    string `json:"image_url"`
}

type Enrollments struct {
	Enrollments []Enrollment `json:"enrollments"`
	Meta        Meta         `json:"meta"`
}

type Enrollment struct {
	UserID          int       `json:"user_id"`
	EnrolledAt      time.Time `json:"enrolled_at"`
	CompletedAt     any       `json:"completed_at"`
	PercentComplete int       `json:"percent_complete"`
	ExpiresAt       any       `json:"expires_at"`
}

type courseData struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Heading  string `json:"heading"`
	Enrolled []User `json:"enrolled"`
}

type Meta struct {
	Total         int `json:"total"`
	Page          int `json:"page"`
	From          int `json:"from"`
	To            int `json:"to"`
	PerPage       int `json:"per_page"`
	NumberOfPages int `json:"number_of_pages"`
}

type Controller struct{}

var client = &http.Client{}
var baseURL = "https://developers.teachable.com/v1"
var allStudents = make(map[int]User)
var allCourses = make(map[int]Course)

func NewController() *Controller {
	return &Controller{}
}

// fetches all users in a school from Teachable API /users endpoint
func (c *Controller) fetchAllUsers() (map[int]User, error) {
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
		// TODO: log error
		return nil, err
	}

	// store all users in-memory by user id
	for _, u := range users.Users {
		var student User
		student.Name = u.Name
		student.Email = u.Email
		allStudents[u.ID] = student
	}
	return allStudents, nil
}

// fetches all courses in a school from Teachable API /courses endpoint
func (controller *Controller) fetchAllCourses() (map[int]Course, error) {
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
		// TODO: log error
		return nil, err
	}

	// store all courses in-memory by course id
	for _, c := range courses.Courses {
		allCourses[c.ID] = c
	}

	return allCourses, nil
}

// fetches all active enrollments by course id from Teachable API /courses/{course_id}/enrollments endpoint
func (controller *Controller) fetchCourseEnrollments(courseID int) ([]User, error) {
	url := fmt.Sprintf("%s/courses/%v/enrollments", baseURL, courseID)
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
		// TODO: log error
		return nil, err
	}

	enrolledStudents := []User{}

	for _, enrolled := range enrollments.Enrollments {
		studentData := allStudents[enrolled.UserID]
		studentData.ID = enrolled.UserID
		enrolledStudents = append(enrolledStudents, studentData)
	}
	return enrolledStudents, nil
}

func (controller *Controller) GetCourseInfo() http.HandlerFunc {
	controller.fetchAllCourses()
	controller.fetchAllUsers()

	return func(w http.ResponseWriter, r *http.Request) {
		// map through each course, and fetch enrolled students by course id.
		// return a list of each published course with course name, heading and enrolled students
		coursesData := controller.fetchEnrollmentsByCourse()
		utils.SendSuccess(w, coursesData)
	}
}

func (controller *Controller) fetchEnrollmentsByCourse() []courseData {
	wg := sync.WaitGroup{}

	var coursesData []courseData

	for _, c := range allCourses {
		if c.IsPublished {
			wg.Add(1)
			go func(course Course) {
				enrollments, err := controller.fetchCourseEnrollments(course.ID)
				if err != nil {
					return
				}
				var data courseData
				data.ID = course.ID
				data.Name = course.Name
				data.Heading = course.Heading
				data.Enrolled = enrollments

				coursesData = append(coursesData, data)
				wg.Done()
			}(c)
		}
	}
	wg.Wait()
	return coursesData
}

func setHeaders(r *http.Request) *http.Request {
	apiKey := os.Getenv("API_KEY")
	r.Header.Add("apiKey", apiKey)
	return r
}
