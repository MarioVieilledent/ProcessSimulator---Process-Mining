# Process Mining Simulation

This program produce an Event Log based on rules.

## Functionality

The program generate randomly an event log in both `JSON` and CSV `format` (JSON is easier to read and CSV can be used for the [open source ProM tool](https://promtools.org/)).
The event log represent **a single process** containing **9 activities** for **100 actors**, here it's students.

- `Activities` are uniquely identified by the property `title`
- `Actors` (here students) are uniquely identified by the property `idRelatedTo`
- `Logs` contains each a timestamp,
    - the property `timestamp` is used for the internal program to handle concurrency
    - the property `date` is a timestamp in ISO format

## Run

1) An executable file for Windows 64bits can be run.

2) Otherwise, requires [golang](https://go.dev/) installed.

Run the program
> `go run .`

Compile the code into an executable
> `go build`

## Rules

### A university receive student who want to register.

There's three roles in the university
- The **Director** who receive the candidatures
- The **Administration** who decide wether to accept or not students in several step of the process
- The **Teachers** who give ranks to student

### This is the 9 tasks/activities for the process:

![Drawn process model before implementing it](./picture_model_dranw.jpg)
*Drawn process model before implementing it*

1) **Receive** *Director*
First a Student send a registration to the university, on the point of view of the university, the director receive the candidature.

2) **Ask for more info** *Administration*
The student matches the criteria so the administration can send him or her an email asking for more information, for example, project works and motivation letter.

3) **Refuse (eligibility)** *Administration*
Maybe some properties related to the student doesn't allow him or her to be enrolled to the university, e.g. he or she is not graduated enough. The administration send to the student an email with the corresponding information.

4) **Ask for register again** *Administration*
Maybe the student didn't give enough information, the administration send to the student an email with the corresponding information.

5) **Refuse (motivation letter)** *Teacher*
Whenever the student that has been asked for give more information sends the information, a teacher read the provided documents and decide the student isn't qualified enough for study in the university. The student is therefore refused.

6) **Rank** *Teacher*
Whenever the student that has been asked for give more information sends the information, a teacher read the provided documents and assign a grade in order to rank the student.

7) **Refuse (grades)** *Administration*
Once the deadline is reached, it's time to choose a certain amount (here in the code a certain proportion) of student regarding to their assigned grade. This task tell the student he or she has been refused.

8) **Accept** *Administration*
Once the deadline is reached, it's time to choose a certain amount (here in the code a certain proportion) of student regarding to their assigned grade. This task tell the student he or she has been accepted and is enrolled in the université.

9) **Too late** *Administration
Once the deadline is reached, administrators tell all students that are not refused or accepted yet that it is now too late to enroll.*

## GoLang

I used Go programming language to implement this program. Go is both easy and simple to use and powerful to handle concurrency.