package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

// Data to work with

var names []string = []string{"James", "John", "Luke", "Jack", "Charles", "Jace", "Chase", "Miles", "Cole", "Max", "Juan", "George", "Blake", "Jayce", "Kai", "Bryce", "King", "Jude", "Grant", "Finn", "Beau", "Mark", "Kyle", "Dean", "Paul", "Zane", "Jax", "Rhett", "Myles", "Brooks", "Sean", "Jase", "Jake", "Knox", "Cash", "Reid", "Chance", "Gage", "Nash", "Lane", "Seth", "Jett", "Troy", "Shane", "Quinn", "Ace", "Colt", "Cruz", "Prince", "Reed", "Frank", "Shawn", "Kash", "Clark", "Jay", "Drew", "Kane", "Wade", "Cade", "Kade", "Zayn", "Hayes", "Bruce", "Tate", "Zayne", "Brock", "Royce", "Scott", "Pierce", "Keith", "Hank", "Rhys"}

var countries []string = []string{"France", "Germany", "Spain", "Italy", "United Kingdom", "Portugal", "Poland", "Sweden", "Norway", "Denmark", "Netherlands", "Switzerland", "Greece", "Russia", "Ukraine", "Slovenia", "Slovakia", "Czechia", "Hungary", "Turkey", "Albania", "Kosovo"}

// Types

type Student struct {
	Id      int    `json:"id"` // Unique id for student
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	State   State  `json:"state"` // State to keep track of the student state
}

// States (for students)
type State string

const WaitingForMoreInfo = "WaitingForMoreInfo" // Student registered
const Refused = "Refused"                       // Reason not relevant for the state
const Ranked = "Ranked"                         // Student has been ranked
const Accepted = "Accepted"                     // The ranking accepted the student

// Roles :  Administration, Director, Teacher
type Role string

const Administration Role = "Administration"
const Director Role = "Director"
const Teacher Role = "Teacher"

// Event = activity
type Event struct {
	Title           string `json:"title" csv:"title"`                    // Title of activity
	Description     string `json:"description"csv:"description"`         // Description of activity
	IdRelatedTo     int    `json:"idRelatedTo"csv:"idRelatedTo"`         // Unique id of user implicated on the event, it is used to group all events to retrieve a single whole process
	NameUser        string `json:"nameUser"csv:"nameUser"`               // Name of the user related (should not be used as unique id for users)
	RoleResponsible Role   `json:"roleResponsible"csv:"roleResponsible"` // Role responsible of doing this activity
	Timestamp       int    `json:"timestamp"csv:"timestamp"`             // Integer as timestamp
	Date            string `json:"date"csv:"date"`                       // Date format for timestamp
}

// Events as constants

// Receive and analyze student registration
const Receive string = "Receive"

// Refuse student's registration for eligibility reasons
const RefuseEligibility string = "Refuse (eligibility)"

// Ask student to register again because of lack of information in registration
const AskRegisterAgain string = "Ask register again"

// Ask student to provide more information (Resume, Transcript of records, Motivation letter)
const AskMoreInfo string = "Ask more info"

// Refuse student application due to teacher's decision over motivation letter
const RefuseMotivationLetter string = "Refuse (motivation letter)"

// Accept student motivation letter and rank him for future and last selection
const Rank string = "Rank"

// Tell student it's too late for application because deadline reached
const TooLate string = "Too late"

// After deadline, refuse ranked student due to final selection (with ranking)
const RefuseGrades string = "Refuse (grades)"

// After deadline, accept ranked student due to final selection (with ranking)
const Accept string = "Accept"

// Variables

var Students []Student = []Student{} // All the students
var EventLog []Event = []Event{}     // The event log is just a list of event
var nbStudents int = 100             // Number of students in applying
var deadLine int = 300               // Timestamp when student registration is over, and ranked bests students are chosen

// Time

var timeModel int = 0 // Time modeled as a simple integer for timestamps

var timeTaskAdministration int = 10_000 // Time in micro second for the administration to process a task
var timeTaskDirector int = 50_000       // Time in micro second for the director to process a task
var TimeRaskTeacher int = 20_000        // Time in micro second for the teacher to process a task

// Channels

var chanRegister chan Student = make(chan Student, 10000) // A Student register for the university
var chanMoreInfo chan Student = make(chan Student, 10000) // A Student gives more information about him/her

var chanAdministration chan Event = make(chan Event, 10000) // A task is sent to the administration service
var chanStudentMail chan Event = make(chan Event, 10000)    // A mail is sent to the student providing information

func main() {
	// Creation of students
	for i := 0; i < nbStudents; i++ {
		createStudent()
	}

	// Start the thread for incrementing time
	go manageTime()

	// Start the thread for university
	go universityThread(chanRegister, chanAdministration)

	// Start the thread for the administration service
	go administrationThread(chanAdministration, chanStudentMail)

	// Starts the thread for students
	go studentsThread(chanRegister, chanStudentMail)

	// Debug (wait some time for all process to end)
	time.Sleep(time.Millisecond * 10_000)

	// Create dates instead of timestamps
	createDates()

	// Write event log in a JSON file
	writeJSON()

	// Write event log in a CSV file
	writeCSV()

	// Check results
	// checkReults()
}

// Debug function to check results at the end of generating process
func checkReults() {
	var list [][]string = make([][]string, len(Students))

	sort.Slice(EventLog, func(i, j int) bool {
		return EventLog[i].Timestamp < EventLog[j].Timestamp
	})

	for _, e := range EventLog {
		list[e.IdRelatedTo] = append(list[e.IdRelatedTo], e.Title)
	}

	for i, process := range list {
		fmt.Print(strconv.Itoa(i) + " => ")
		for _, task := range process {
			fmt.Print(task, ", ")
		}
		fmt.Println()
	}

}

// Behavior of the university
func universityThread(chanRegister chan Student, chanAdministration chan Event) {
	for {
		select {
		case s1 := <-chanRegister:
			processRegistration(s1)
		case s2 := <-chanMoreInfo:
			processMoreInfo(s2)
		}
	}
}

// Director's role
func processRegistration(s Student) {
	time.Sleep(time.Duration(timeTaskDirector) * time.Microsecond)
	if timeModel < deadLine {
		task := Event{
			Title:           Receive,
			Description:     "Receive register request from " + s.Name + ", " + strconv.Itoa(s.Age) + " years old, from " + s.Country + ", the request will be analyzed soon enough.",
			IdRelatedTo:     s.Id,
			NameUser:        s.Name,
			RoleResponsible: Director,
			Timestamp:       timeModel,
			Date:            "",
		}
		log(task)
		chanAdministration <- task
	}
}

// Teacher's role
func processMoreInfo(s Student) {
	time.Sleep(time.Duration(TimeRaskTeacher) * time.Microsecond)
	task := Event{
		Title:           "",
		Description:     "",
		IdRelatedTo:     s.Id,
		NameUser:        s.Name,
		RoleResponsible: Teacher,
		Timestamp:       timeModel,
		Date:            "",
	}
	// One chance over 5 to be refused, otherwise the student is getting a rank
	if rand.Intn(4) == 0 {
		Students[task.IdRelatedTo].State = Refused
		task.Title = RefuseMotivationLetter
		task.Description = "We are deeply sorry to announce you that out university judged your grades too low regarding to information you provided. We are glad you took the application in consideration and wish you luck for your future studies."
	} else {
		Students[task.IdRelatedTo].State = Ranked
		task.Title = Rank
		task.Description = "We are happy to announce you that you have been selected for ranking. Because we receive a lot of requests from students, we have to rank all students regarding to their grades. You will receive a mail telling you if you have been accepted by the ranking."
	}
	log(task)
	chanAdministration <- task
}

// Behavior of administration
func administrationThread(chanAdministration, chanStudentMail chan Event) {
	for {
		adminTask := <-chanAdministration
		switch adminTask.Title {
		case Receive:
			receive(adminTask, chanStudentMail)
		}
	}
}

// Administration has to check the registration of student
func receive(adminTask Event, chanStudentMail chan Event) {
	time.Sleep(time.Duration(timeTaskAdministration) * time.Microsecond)
	responseTask := Event{
		Title:           "",
		Description:     "",
		IdRelatedTo:     adminTask.IdRelatedTo,
		NameUser:        adminTask.NameUser,
		RoleResponsible: Administration,
		Timestamp:       timeModel,
		Date:            "",
	}
	if Students[adminTask.IdRelatedTo].Age < 18 {
		Students[adminTask.IdRelatedTo].State = Refused
		responseTask.Title = RefuseEligibility
		responseTask.Description = "The university does not accept student under 18 years old."
		chanStudentMail <- responseTask
	} else if Students[adminTask.IdRelatedTo].Country == "France" {
		Students[adminTask.IdRelatedTo].State = Refused
		responseTask.Title = RefuseEligibility
		responseTask.Description = "The university does not accept student from France due to previous bad experience with them."
		chanStudentMail <- responseTask
	} else if rand.Intn(4) == 0 { // Else one time over 5, we assumes student didn't gave enough information for registration
		responseTask.Title = AskRegisterAgain
		responseTask.Description = "The register you provided student is not filled correctly, please register again with correct information."
		chanStudentMail <- responseTask
	} else {
		responseTask.Title = AskMoreInfo
		responseTask.Description = "Your register has been accepted by Administration team. Now you need to provide us more information (transcript of records and motivation letter). Please send us the document before the deadline."
		chanStudentMail <- responseTask
	}
	log(responseTask)
}

// Behavior of a student
func studentsThread(chanRegister chan Student, chanStudentMail chan Event) {
	// For all students, send a registration
	for _, s := range Students {
		go register(s, chanRegister)
	}

	// This loop receive all students tasks
	for {
		studentTask := <-chanStudentMail
		switch studentTask.Title {
		case AskMoreInfo:
			askMoreInfo(studentTask)
		case AskRegisterAgain:
			askRegisterAgain(studentTask)
		case RefuseEligibility:
			refuseEligibility(studentTask)
		}
	}
}

// A student register (with a random delay)
func register(s Student, chanRegister chan Student) {
	// Wait a random amount of time
	time.Sleep(time.Duration(randomTimestamp(0, 50)) * time.Millisecond)
	chanRegister <- s
}

// A student is refused by university (the student can't do anything else)
func refuseEligibility(studentTask Event) {}

// A student omit to mention his/her country, he/she has to send registration again
func askRegisterAgain(studentTask Event) {
	register(Students[studentTask.IdRelatedTo], chanRegister)
}

// A Student has been accepted and asked to send more info to the university in order to be ranked
func askMoreInfo(studentTask Event) {
	// Wait a random amount of time
	time.Sleep(time.Duration(randomTimestamp(0, 50)) * time.Millisecond)
	chanMoreInfo <- Students[studentTask.IdRelatedTo]
}

// Add 1 to time variable each specific amount of time (I think millisecond to process the whole program quickly)
func manageTime() {
	for {
		time.Sleep(1 * time.Millisecond)
		timeModel += 1
		// fmt.Println("Time =", timeModel)

		if timeModel == deadLine {
			go deadlineReached()
		}
	}
}

func deadlineReached() {
	// Selection of best ranked student
	go selectBestStudents()

	// For all student who is refused or who didn't send their info (or if the university didn't had enough time to check more info), the university tell him/her it's now too late
	for _, s := range Students {
		if s.State == WaitingForMoreInfo {
			time.Sleep(time.Duration(timeTaskAdministration) * time.Microsecond)
			s.State = Refused
			tooLateTask := Event{
				Title:           TooLate,
				Description:     "We are sorry to announce you that it is now too late for applying to the university. Thank you for your interest on working with us.",
				IdRelatedTo:     s.Id,
				NameUser:        s.Name,
				RoleResponsible: Administration,
				Timestamp:       timeModel,
				Date:            "",
			}
			log(tooLateTask)
		}
	}
}

// When deadline is reached, select 75% of ranked student (simulation of choosing best student by randomly choosing ranked students)
func selectBestStudents() {
	for _, s := range Students {
		if s.State == Ranked {
			time.Sleep(time.Duration(timeTaskAdministration) * time.Microsecond)
			task := Event{
				Title:           "",
				Description:     "",
				IdRelatedTo:     s.Id,
				NameUser:        s.Name,
				RoleResponsible: Administration,
				Timestamp:       timeModel,
				Date:            "",
			}
			if rand.Intn(3) == 0 { // 25 % chance to be refused because of bad grades
				s.State = Refused
				task.Title = RefuseGrades
				task.Description = "We are deeply sorry to announce you that we can not keep you in our University because of too much application. We are very glad of your interest in our university and wish you best luck for your future studies."
			} else { // 75 % chance to be accepted
				s.State = Accepted
				task.Title = Accept
				task.Description = "We are deeply sorry to announce you that we can not keep you in our University because of too much application. We are very glad of your interest in our university and wish you best luck for your future studies."
			}
			log(task)
		}
	}
}

// Create a random student
func createStudent() {
	Students = append(Students, Student{len(Students), randomName(), randomAge(), randomCountry(), WaitingForMoreInfo})
}

// Get a random name
func randomName() string {
	return names[rand.Intn(len(names))]
}

// Get a random integer age between min and max
func randomAge() int {
	max := 33
	min := 15
	return rand.Intn(max-min) + min
}

// Get a random country
func randomCountry() string {
	return countries[rand.Intn(len(countries))]
}

// Get a random timestamp between min and max
func randomTimestamp(min, max int) int {
	return rand.Intn(max-min) + min
}

// Append event to the event log
func log(task Event) {
	EventLog = append(EventLog, task)
	fmt.Println("Event logged => " + task.Title + " [" + strconv.Itoa(task.IdRelatedTo) + "]")
}

// Create dates instead of timestamps
func createDates() {
	fmt.Println("Debug")
	for i, event := range EventLog {
		minuts := event.Timestamp / 60
		seconds := event.Timestamp % 60
		var dateStr = "2022-12-14T05:" + conv(minuts) + ":" + conv(seconds) + "Z"
		EventLog[i].Date = dateStr
	}
}

// Convert int to string minut / second format with extra 0 at the start if needed
func conv(val int) string {
	if val < 10 {
		return "0" + strconv.Itoa(val)
	} else {
		return strconv.Itoa(val)
	}
}

// Write eventLog file in json
func writeJSON() {
	result, err := json.MarshalIndent(EventLog, "", "\t")
	if err != nil {
		fmt.Println("Error stringifying JSON :", err)
	}

	// d1 := []byte(string(result))

	f, err := os.Create("./eventLog.json")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(string(result))

	if err2 != nil {
		panic(err2)
	}

	fmt.Println("eventLog JSON file successfully written.")
}

// Write eventLog file in csv
func writeCSV() {
	// to download file inside downloads folder
	clientsFile, err := os.OpenFile("./eventLog.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	err = gocsv.MarshalFile(&EventLog, clientsFile) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}
