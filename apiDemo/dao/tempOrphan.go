package dao

import (
	"github.com/dileepaj/tracified-gateway/model"
)

func (cd *Connection) InsertTempOrphan(tempOrphan model.TransactionCollectionBody) (string,error, int) {
	return Create(tempOrphan, "TempOrphan")

}
