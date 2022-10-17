package authentication

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type PermissionStatus struct {
	Status             bool
	TenantId           string
	UserId             string
	IsSubscriptionPaid bool
}

/*
 * The HasPermission function will return a boolean value which will be used to check if the user has the required access claim or not
 */

func HasPermission(reqToken string) PermissionStatus {
	var ps PermissionStatus

	if len(reqToken) > 0 {
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("bsof2sJXPp0T5G38L6RKq21mqayXyr4u"), nil //todo move this to env file
		})
		if err != nil {
			if err.Error() == jwt.ErrSignatureInvalid.Error() {
				logrus.Println(err.Error())
				return ps
			}
			logrus.Println(err.Error())
			return ps
		}

		for key, val := range claims {
			if key == "userID" {
				ps.UserId = fmt.Sprintf("%v", val)
			}
			if key == "tenantID" {
				ps.TenantId = fmt.Sprintf("%v", val)
			}
			if key == "permissions" {
				v, ok := val.(map[string]interface{})["0"]
				if !ok {
					logrus.Println("Permissions not found")
				}
				if v != nil {
					switch reflect.TypeOf(v).Kind() {
					case reflect.Slice:
						s := reflect.ValueOf(v)
						for i := 0; i < s.Len(); i++ {
							if s.Index(i).Interface().(string) == "98" {
								ps.Status = true
							}
						}
					}
				} else {
					logrus.Println("Permissions not found")
					ps.Status = false
				}
			}
			if key == "isSubscriptionPaid" {
				ps.IsSubscriptionPaid, _ = strconv.ParseBool(fmt.Sprintf("%v", val))
			}
		}
	} else {
		logrus.Println("Bearer token not found")
		return ps
	}
	return ps
}
