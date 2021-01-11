package dao

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"gopkg.in/mgo.v2"
)

/*Connection The Mgo Connection
@author - Azeem Ashraf
*/
type Connection struct {
}

func (cd *Connection) connect()(*mgo.Session,error) {
	return commons.GetMongoSession()
}
