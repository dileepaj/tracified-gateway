package model
import "time"

type StellarTransaction struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Account struct {
			Href string `json:"href"`
		} `json:"account"`
		Ledger struct {
			Href string `json:"href"`
		} `json:"ledger"`
		Operations struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"operations"`
		Effects struct {
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"effects"`
		Precedes struct {
			Href string `json:"href"`
		} `json:"precedes"`
		Succeeds struct {
			Href string `json:"href"`
		} `json:"succeeds"`
		Transaction struct {
			Href string `json:"href"`
		} `json:"transaction"`
	} `json:"_links"`
	ID                    string    `json:"id"`
	PagingToken           string    `json:"paging_token"`
	Successful            bool      `json:"successful"`
	Hash                  string    `json:"hash"`
	Ledger                int       `json:"ledger"`
	CreatedAt             time.Time `json:"created_at"`
	SourceAccount         string    `json:"source_account"`
	SourceAccountSequence string    `json:"source_account_sequence"`
	FeeAccount            string    `json:"fee_account"`
	FeeCharged            string    `json:"fee_charged"`
	MaxFee                string    `json:"max_fee"`
	OperationCount        int       `json:"operation_count"`
	EnvelopeXdr           string    `json:"envelope_xdr"`
	ResultXdr             string    `json:"result_xdr"`
	ResultMetaXdr         string    `json:"result_meta_xdr"`
	FeeMetaXdr            string    `json:"fee_meta_xdr"`
	MemoType              string    `json:"memo_type"`
	Signatures            []string  `json:"signatures"`
}

type StellarOperations struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		Prev struct {
			Href string `json:"href"`
		} `json:"prev"`
	} `json:"_links"`
	Embedded struct {
		Records []struct {
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				Transaction struct {
					Href string `json:"href"`
				} `json:"transaction"`
				Effects struct {
					Href string `json:"href"`
				} `json:"effects"`
				Succeeds struct {
					Href string `json:"href"`
				} `json:"succeeds"`
				Precedes struct {
					Href string `json:"href"`
				} `json:"precedes"`
			} `json:"_links"`
			ID                    string    `json:"id"`
			PagingToken           string    `json:"paging_token"`
			TransactionSuccessful bool      `json:"transaction_successful"`
			SourceAccount         string    `json:"source_account"`
			Type                  string    `json:"type"`
			TypeI                 int       `json:"type_i"`
			CreatedAt             time.Time `json:"created_at"`
			TransactionHash       string    `json:"transaction_hash"`
			Name                  string    `json:"name"`
			Value                 string    `json:"value"`
		} `json:"records"`
	} `json:"_embedded"`
}

type ManageData struct {
	Name           string
	Value          string
	Source_account string
	Asset_code     string
	Amount         string
	To             string
	From           string
}