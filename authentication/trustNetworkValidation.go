package authentication

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/sirupsen/logrus"
)

/*
 * This function is used to validate the expert public key against the trust network
 * If the expert is not in the trust network then it will return the error
 * else it will return nil

 * @param expertPK - expert public key

 * A collection is used to store the trust network key map
 */

func ValidateAgainstTrustNetwork(expertPK string) error {
	object := dao.Connection{}
	TrustNetworkKeyMap, err := object.GetTrustNetworkKeyMap(expertPK).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Error("Error while getting the trust network key map : ", err)
		return errors.New("Error while getting the trust network key map")
	}
	if TrustNetworkKeyMap == nil {
		return errors.New("Expert is not in the trust network")
	} else {
		return nil
	}
}
