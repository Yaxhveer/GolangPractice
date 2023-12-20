package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Model for course
type Course struct {
	CourseId string `json:"courseid"`
	CourseName string `json:"coursename"`
	CoursePrice int `json:"price"`
	Author *Author `json:"author"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website string `json:"website"`
}

var courses []Course;

// Middleware
func (c *Course) IsEmpty() bool {
	return c.CourseName == "";
}

func main()  {
	fmt.Println("Yooo");

	r := mux.NewRouter();

	author := Author{Fullname: "Yash", Website: "yash@gmail.com"}
	courses = append(courses, Course{CourseId: "1", CourseName: "Hell", CoursePrice: 0, Author: &author})
	courses = append(courses, Course{CourseId: "2", CourseName: "Hell2", CoursePrice: 20, Author: &author})

	r.HandleFunc("/", serveHome).Methods("GET");
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getCourse).Methods("GET")
	r.HandleFunc("/course", createCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateCourse).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteCourse).Methods("DELETE")

	http.Handle("/", r);
	log.Fatal(http.ListenAndServe(":3000", r));
};

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home Page")
	w.Write([]byte("<h1>Home</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get All Course")
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Course");
	w.Header().Add("content-type", "application/json");

	// get id
	params := mux.Vars(r);

	for _, course := range courses {
		if (course.CourseId == params["id"]){
			json.NewEncoder(w).Encode(course);
			return;
		}
	}
	json.NewEncoder(w).Encode("No Course found.");
}

func createCourse (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add one Course");
	w.Header().Add("content-type", "application/json");
	fmt.Println(r.Body)

	if r.Body == nil {
		json.NewEncoder(w).Encode("Send some data")
		return;
	}

	var course Course;
	_ = json.NewDecoder(r.Body).Decode(&course);

	fmt.Println(course);

	if course.IsEmpty() {
		json.NewEncoder(w).Encode("JSON is empty.")
		return;
	}

	// generating new id
	course.CourseId = strconv.Itoa(rand.Intn(100));
	courses = append(courses, course);
	json.NewEncoder(w).Encode("Task completed")
}

func updateCourse (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Course");
	w.Header().Add("content-type", "application/json");

	params := mux.Vars(r)

	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)

			var cc Course;
			_ = json.NewDecoder(r.Body).Decode(&cc)

			cc.CourseId = params["id"]
			courses = append(courses, cc)
			json.NewEncoder(w).Encode("Task completed")

			return
		}
	}
}

func deleteCourse(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Delete a Course");
	w.Header().Add("content-type", "application/json");

	params := mux.Vars(r);

	for index, course := range courses {
		if course.CourseId == params["id"]{
			courses = append(courses[:index], courses[index+1:]...);
			break;
		}
	}
	json.NewEncoder(w).Encode("Deleted")
}