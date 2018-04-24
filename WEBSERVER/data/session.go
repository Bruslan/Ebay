package data

import (
	"github.com/gocql/gocql"
	"time"
)

type Session struct {
	Uuid      string
	UserId    gocql.UUID
	CreatedAt time.Time
	Device    string
	Active    bool
}

// Check if session is valid in the database ***
func (sess *Session) SessValid() bool {

	// database session
	dbs, err := db_session()
	if err != nil {
		return false
	}
	defer dbs.Close()
	var device string
	if err = dbs.Query(`SELECT active, device FROM user_keyspace.sessions WHERE uuid=?`, sess.Uuid).Scan(
		&sess.Active,
		&device); err != nil {
		return false
	} else {
		if sess.Active && (sess.Device == device) {
			return true
		} else {
			return false
		}
	}
}

// checks if user id session exists with same device but inactive
// if such a session exists delete
func (sess *Session) InactiveExists() error {

	// database session
	dbs, err := db_session()
	if err != nil {
		return err
	}
	defer dbs.Close()

	var dev, sessid string
	it := dbs.Query(`SELECT uuid, device FROM user_keyspace.sessions WHERE user_id=?`, sess.UserId).Iter()
	defer it.Close()
	if it.NumRows() < 1 {
		return err
	} else if it.NumRows() > 4 {
		for it.Scan(&sessid, &dev) {
			err = dbs.Query(`DELETE FROM user_keyspace.sessions WHERE uuid=? IF EXISTS`, sessid).Exec()
			if err != nil {
				return err
			}
		}
	} else {
		for it.Scan(&sessid, &dev) {
			if dev == sess.Device {
				err = dbs.Query(`DELETE FROM user_keyspace.sessions WHERE uuid=? IF EXISTS`, sessid).Exec()
			}
		}
	}
	return err
}

// ******************************
// modify this function and delete only session which device information
// Delete session from database
func (sess *Session) Delete() (err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	return dbs.Query(`DELETE FROM user_keyspace.sessions WHERE uuid=? IF EXISTS`, sess.Uuid).Exec()
}

func (sess *Session) SetInactive() (err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	return dbs.Query(`UPDATE user_keyspace.sessions SET active=? WHERE uuid=? IF EXISTS`, false, sess.Uuid).Exec()
}

// Get the user from the session
func (session *Session) User() (user User, err error) {

	// database session
	dbs, err := db_session()
	if err != nil {
		return
	}
	defer dbs.Close()
	user = User{}
	err = dbs.Query(`SELECT user_id FROM user_keyspace.sessions WHERE uuid=?`, session.Uuid).Scan(
		&user.Id)
	return
}
