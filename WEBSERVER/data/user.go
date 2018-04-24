package data

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"strings"
	"time"
)

type User struct {
	Id        gocql.UUID
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Country   string
	Password  string
	Birthday  time.Time
	CreatedAt time.Time
}

var cluster *gocql.ClusterConfig

// init cassandra configuration
func init() {
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "user_keyspace"
	cluster.Consistency = gocql.One
}

// create cassandra session
func db_session() (dbs *gocql.Session, err error) {
	dbs, err = cluster.CreateSession()
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

// Create a new session for an existing user
func (user *User) CreateSession(device string) (ses Session, err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	// uuid creation
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return
	}
	ses = Session{Uuid: uuid.String(), UserId: user.Id, CreatedAt: time.Now(), Device: device, Active: true}
	// insert session into database
	if err = dbs.Query(`INSERT INTO user_keyspace.sessions (uuid, user_id, created_at, device, active) VALUES (?, ?, ?, ?, ?) IF NOT EXISTS`,
		ses.Uuid,
		ses.UserId,
		ses.CreatedAt,
		ses.Device,
		ses.Active).Exec(); err != nil {
		return
	}
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	// session = Session{}
	// err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
	// 	Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error, stmt string) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	// create new uuid
	uuid, err := gocql.RandomUUID()
	if err != nil {
		stmt = "We apologize. Something went wrong, please try again."
		return
	}
	// check if username, email exists
	i1 := dbs.Query(`SELECT * FROM user_keyspace.users WHERE username=?`, user.UserName).Iter()
	i2 := dbs.Query(`SELECT * FROM user_keyspace.users WHERE email=?`, user.Email).Iter()
	defer i1.Close()
	defer i2.Close()

	if i1.NumRows() > 0 {
		err = errors.New("User already exists")
		stmt = "This Username already exists, please try another one."
		return
	} else if i2.NumRows() > 0 {
		err = errors.New("Email already exists")
		stmt = "This Email already exists, please try another one."
		return
	} else {
		// insert user into database
		if err := dbs.Query(`INSERT INTO user_keyspace.users (id, username, first_name, last_name, email, pass, country, birthday, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) IF NOT EXISTS`,
			uuid,
			user.UserName,
			user.FirstName,
			user.LastName,
			user.Email,
			Encrypt(user.Password),
			user.Country,
			user.Birthday,
			time.Now()).Exec(); err != nil {
			log.Println("Create user Error:", err)
		} else {
			fmt.Println("inserted user into database")
		}
	}
	return
}

// Delete user from database
func (user *User) Delete() (err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	return dbs.Query(`DELETE FROM user_keyspace.users WHERE id=? IF EXISTS`, user.Id).Exec()
}

// Delete all sessions of user
func (user *User) DeleteSessions() (err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	var uuid string
	// get all session uuids
	it := dbs.Query(`SELECT uuid FROM user_keyspace.sessions WHERE user_id=?`, user.Id.String()).Iter()
	for it.Scan(&uuid) {
		err = dbs.Query(`DELETE FROM user_keyspace.sessions WHERE uuid=? IF EXISTS`, uuid).Exec()
	}
	defer it.Close()
	return
}

// // Update user information in the database
// func (user *User) Update() (err error) {
// 	statement := "update users set name = $2, email = $3 where id = $1"
// 	stmt, err := Db.Prepare(statement)
// 	if err != nil {
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(user.Id, user.Name, user.Email)
// 	return
// }

// // Get all users in the database and returns it
// func Users() (users []User, err error) {
// 	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
// 	if err != nil {
// 		return
// 	}
// 	for rows.Next() {
// 		user := User{}
// 		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
// 			return
// 		}
// 		users = append(users, user)
// 	}
// 	rows.Close()
// 	return
// }

// Get a single user given the email
func GetUserEmail(nameemail string) (user User, err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()

	user = User{}
	if strings.Contains(nameemail, "@") {
		err = dbs.Query(`SELECT id, pass FROM user_keyspace.users WHERE email=?`, nameemail).Scan(
			&user.Id,
			&user.Password)
	} else {
		err = dbs.Query(`SELECT id, pass FROM user_keyspace.users WHERE username=?`, nameemail).Scan(
			&user.Id,
			&user.Password)
	}
	return
}

// Get a single user given the UUID
func GetUserByUuid(uuid string) (user User, err error) {
	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	user = User{}

	// convert to database uuid type
	go_uuid, err := gocql.ParseUUID(uuid)
	if err != nil {
		return
	}

	err = dbs.Query(`SELECT id, email, name, pass FROM user_keyspace.users WHERE id=?`, go_uuid).Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.Password)
	return
}
