package data

import (
	"fmt"
	"github.com/gocql/gocql"
	"testing"
	"time"
)

// test data
var user = User{
	UserName: "freezy",
	Email:    "jansifreezy@koko.com",
	Password: "peter_pass",
}

// func setup() {
// 	ThreadDeleteAll()
// 	SessionDeleteAll()
// 	UserDeleteAll()
// }

func TestInit(t *testing.T) {
	fmt.Println("Testing Print")
	// connStr := "user=postjansen password=postjansen93 dbname=jansendb sslmode=disable"
	// _, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	t.Errorf("Cannot open Database connection, err: %v", err)
	// }

	uuid1, err := gocql.RandomUUID()
	if err != nil {
		fmt.Println("uuid1 error")
	}
	fmt.Println("UUID=", uuid1)
	fmt.Println(uuid1, "amanakoidum")
	fmt.Println(uuid1.String())

	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "user_keyspace"
	cluster.Consistency = gocql.One
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "ianzndb",
		Password: "Lov3toN8t",
	}
	session, err := cluster.CreateSession()
	fmt.Println("Session error:", err)

	defer session.Close()

	fmt.Println(user.Email)
	fmt.Println(user.Name)

	// create new uuid
	uuid, err := gocql.RandomUUID()
	if err != nil {
		fmt.Println("Could not create new User ID")
	}
	// check if username exists
	iter := session.Query(`SELECT * FROM user_keyspace.users WHERE name=?`, user.Name).Iter()
	fmt.Println("iterator length:", iter.NumRows(), iter)
	// check if email exists:
	iter2 := session.Query(`SELECT * FROM user_keyspace.users WHERE email=?`, user.Email).Iter()
	fmt.Println("iterator length:", iter2.NumRows(), iter2)

	if iter.NumRows() == 0 && iter2.NumRows() == 0 {
		// INSERT INTO user_keyspace.users (id, created_at, email, name, password, uuid) VALUES (1, toTimestamp(now()), 'mongo@holiday.com', 'mr.bongo', 'lalala', now());
		// insert a tweet
		if err := session.Query(`INSERT INTO user_keyspace.users (id, city, created_at, email, name, pass) VALUES (?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
			uuid, "rom", time.Now(), "jansifreezy@koko.com", "freezy", "freezy_pw").Exec(); err != nil {
			fmt.Println("Inserted values and thats the err:", err)
		}
	}

	var id gocql.UUID
	var emailo, namee string
	var created time.Time

	email := ""

	// first query
	if err := session.Query(`SELECT id, email, created_at FROM user_keyspace.users`).Scan(&id, &email, &created); err != nil {
		fmt.Println("First query:", err)
	}
	fmt.Println("First query, im proud :)", id, email, created)

	// list all tweets
	iter = session.Query(`SELECT id, email, name FROM user_keyspace.users`).Iter()
	for iter.Scan(&id, &emailo, &namee) {
		fmt.Println("Iteration query boy", id, emailo, namee)
	}
	if err := iter.Close(); err != nil {
		fmt.Println(err)
	}

}

// func TestQuery(t *testing.T) {
// 	Db := *sql.DBuser_keyspace.users

// 	connStr := "user=postjansen password=postjansen93 dbname=jansendb sslmode=disable"
// 	Db, err := sql.Open("postgres", connStr)

// }

func TestCreateUUID(t *testing.T) {

}

// import (
// 	"database/sql"
// 	"testing"
// )

// // test data
// var users = []User{
// 	{
// 		Name:     "Peter Jones",
// 		Email:    "peter@gmail.com",
// 		Password: "peter_pass",
// 	},
// 	{
// 		Name:     "John Smith",
// 		Email:    "john@gmail.com",
// 		Password: "john_pass",
// 	},
// }

// func setup() {
// 	ThreadDeleteAll()
// 	SessionDeleteAll()
// 	UserDeleteAll()
// }

// func Test_UserCreate(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	if users[0].Id == 0 {
// 		t.Errorf("No id or created_at in user")
// 	}
// 	u, err := UserByEmail(users[0].Email)
// 	if err != nil {
// 		t.Error(err, "User not created.")
// 	}
// 	if users[0].Email != u.Email {
// 		t.Errorf("User retrieved is not the same as the one created.")
// 	}
// }

// func Test_UserDelete(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	if err := users[0].Delete(); err != nil {
// 		t.Error(err, "- Cannot delete user")
// 	}
// 	_, err := UserByEmail(users[0].Email)
// 	if err != sql.ErrNoRows {
// 		t.Error(err, "- User not deleted.")
// 	}
// }

// func Test_UserUpdate(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	users[0].Name = "Random User"
// 	if err := users[0].Update(); err != nil {
// 		t.Error(err, "- Cannot update user")
// 	}
// 	u, err := UserByEmail(users[0].Email)
// 	if err != nil {
// 		t.Error(err, "- Cannot get user")
// 	}
// 	if u.Name != "Random User" {
// 		t.Error(err, "- User not updated")
// 	}
// }

// func Test_UserByUUID(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	u, err := UserByUUID(users[0].Uuid)
// 	if err != nil {
// 		t.Error(err, "User not created.")
// 	}
// 	if users[0].Email != u.Email {
// 		t.Errorf("User retrieved is not the same as the one created.")
// 	}
// }

// func Test_Users(t *testing.T) {
// 	setup()
// 	for _, user := range users {
// 		if err := user.Create(); err != nil {
// 			t.Error(err, "Cannot create user.")
// 		}
// 	}
// 	u, err := Users()
// 	if err != nil {
// 		t.Error(err, "Cannot retrieve users.")
// 	}
// 	if len(u) != 2 {
// 		t.Error(err, "Wrong number of users retrieved")
// 	}
// 	if u[0].Email != users[0].Email {
// 		t.Error(u[0], users[0], "Wrong user retrieved")
// 	}
// }

// func Test_CreateSession(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	session, err := users[0].CreateSession()
// 	if err != nil {
// 		t.Error(err, "Cannot create session")
// 	}
// 	if session.UserId != users[0].Id {
// 		t.Error("User not linked with session")
// 	}
// }

// func Test_GetSession(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	session, err := users[0].CreateSession()
// 	if err != nil {
// 		t.Error(err, "Cannot create session")
// 	}

// 	s, err := users[0].Session()
// 	if err != nil {
// 		t.Error(err, "Cannot get session")
// 	}
// 	if s.Id == 0 {
// 		t.Error("No session retrieved")
// 	}
// 	if s.Id != session.Id {
// 		t.Error("Different session retrieved")
// 	}
// }

// func Test_checkValidSession(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	session, err := users[0].CreateSession()
// 	if err != nil {
// 		t.Error(err, "Cannot create session")
// 	}

// 	uuid := session.Uuid

// 	s := Session{Uuid: uuid}
// 	valid, err := s.Check()
// 	if err != nil {
// 		t.Error(err, "Cannot check session")
// 	}
// 	if valid != true {
// 		t.Error(err, "Session is not valid")
// 	}

// }

// func Test_checkInvalidSession(t *testing.T) {
// 	setup()
// 	s := Session{Uuid: "123"}
// 	valid, err := s.Check()
// 	if err == nil {
// 		t.Error(err, "Session is not valid but is validated")
// 	}
// 	if valid == true {
// 		t.Error(err, "Session is valid")
// 	}

// }

// func Test_DeleteSession(t *testing.T) {
// 	setup()
// 	if err := users[0].Create(); err != nil {
// 		t.Error(err, "Cannot create user.")
// 	}
// 	session, err := users[0].CreateSession()
// 	if err != nil {
// 		t.Error(err, "Cannot create session")
// 	}

// 	err = session.DeleteByUUID()
// 	if err != nil {
// 		t.Error(err, "Cannot delete session")
// 	}
// 	s := Session{Uuid: session.Uuid}
// 	valid, err := s.Check()
// 	if err == nil {
// 		t.Error(err, "Session is valid even though deleted")
// 	}
// 	if valid == true {
// 		t.Error(err, "Session is not deleted")
// 	}
// }
