package dao

import (
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/model"
)

/*InsertCoc Insert a single COC Object to COCCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertCoc(Coc model.COCCollectionBody) error {
	fmt.Println("--------------------------- InsertCoc ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("COC")
	err = c.Insert(Coc)
	if err != nil {
		fmt.Println("Error while inserting to COC " + err.Error())
	}
	return err
}

/*InsertTransaction Insert a single Transaction Object to TransactionCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertTransaction(Coc model.TransactionCollectionBody) error {
	log.Println("--------------------------- InsertTransaction ------------------------")
	session, err := cd.connect()
	if err != nil {
		log.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Transactions")
	err = c.Insert(Coc)
	if err != nil {
		log.Println("Error while inserting to Transactions " + err.Error())
	}
	return err
}

/*InsertSpecialToTempOrphan Insert a single Transaction Object to TempOrphan in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertSpecialToTempOrphan(Coc model.TransactionCollectionBody) error {
	fmt.Println("--------------------------- InsertSpecialToTempOrphan ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("TempOrphan")
	err = c.Insert(Coc)
	if err != nil {
		fmt.Println("Error while inserting to TempOrphan " + err.Error())
	}
	return err
}

/*InsertToOrphan Insert a single Transaction Object to OrphanCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertToOrphan(Coc model.TransactionCollectionBody) error {
	fmt.Println("--------------------------- InsertToOrphan ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Orphan")
	err = c.Insert(Coc)
	if err != nil {
		fmt.Println("Error while inserting to Orphan " + err.Error())
	}
	return err
}

/*InsertProfile Insert a single Profile Object to ProfileCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertProfile(Coc model.ProfileCollectionBody) error {
	fmt.Println("--------------------------- InsertProfile ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Profiles")
	err = c.Insert(Coc)
	if err != nil {
		fmt.Println("Error while inserting to Profiles " + err.Error())
	}
	return err
}

/*InsertCertificate Insert a single Certificate Object to CertificateCollection in DB
@author - Azeem Ashraf
*/
func (cd *Connection) InsertCertificate(Cert model.CertificateCollectionBody) error {
	fmt.Println("--------------------------- InsertCertificate ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Certificates")
	err = c.Insert(Cert)
	if err != nil {
		fmt.Println("Error while inserting to Certificates " + err.Error())
	}
	return err
}

func (cd *Connection) InsertArtifact(artifacts model.ArtifactTransaction) error {
	fmt.Println("--------------------------- InsertArtifact ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Artifacts")
	err = c.Insert(artifacts)
	if err != nil {
		fmt.Println("Error while inserting to Artifacts " + err.Error())
	}
	return err
}

func (cd *Connection) InsertOrganization(Org model.TestimonialOrganization) error {
	fmt.Println("--------------------------- InsertOrganization ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Organizations")
	err = c.Insert(Org)
	if err != nil {
		fmt.Println("Error while inserting to organizations " + err.Error())
	}
	return err
}

func (cd *Connection) InsertTestimonial(Tes model.Testimonial) error {
	fmt.Println("--------------------------- InsertOrganization ------------------------")
	session, err := cd.connect()
	if err != nil {
		fmt.Println("Error while getting session " + err.Error())
	}
	defer session.Close()
	c := session.DB(dbName).C("Testimonials")
	err = c.Insert(Tes)
	if err != nil {
		fmt.Println("Error while inserting to organizations " + err.Error())
	}
	return err
}
