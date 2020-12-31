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

  

  
// 	session, err := mgo.Dial("mongodb://Zeemzo:abcd1234@ds143953.mlab.com:43953/tracified-gateway")
	
	//mongo connection to 99xnsbm Mlab Account
	session, err := mgo.Dial("mongodb://gateway-user:GW%40pass123@db.tracified.com:27017/tracified-gateway")

  
	if err != nil {
		fmt.Println(err)
	}
	return session,err

}
