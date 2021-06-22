package adminDAO

type FOPK struct {
	Accounts []Account `bson:"accounts,omitempty"`
}

type Account struct {
	Pk string `bson:"pk,omitempty"`
}
