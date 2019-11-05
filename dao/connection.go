package dao

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

/*Connection The Mgo Connection
@author - Azeem Ashraf
*/
type Connection struct {
}

func (cd *Connection) connect()(*mgo.Session,error) {
	//mongo connection to Zeemzo Mlab Account
	// session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds143953.mlab.com:43953/tracified-gateway")
	
	//mongo connection to 99xnsbm Mlab Account production
	// session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds157654.mlab.com:57654/tracified-gateway")

	//mongo connection to 99xnsbm Mlab Account staging
	session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds237588.mlab.com:37588/tracified-gateway-staging")
	if err != nil {
		fmt.Println(err)
	}
	return session,err

}
