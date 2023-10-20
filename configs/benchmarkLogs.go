package configs

type benchmarkLogsTag struct {
	TDP_REQUEST              string
	TDP_REQUEST_GENESIS      string
	TDP_REQUEST_NORMAL       string
	TDP_REQUEST_SPLIT_PARENT string
	TDP_REQUEST_SPLIT_CHILD  string
	TDP_REQUEST_MERGE_7      string
	TDP_REQUEST_MERGE_8      string
	TDP_REQUEST_TRANSFER     string
}

type benchmarkLogsActions struct{}

var BenchmarkLogsTag = benchmarkLogsTag{
	"tdp-request",
	"tdp-request-genesis",
	"tdp-request-normal",
	"tdp-request-split-parent",
	"tdp-request-split-child",
	"tdp-request-merge-7",
	"tdp-request-merge-8",
	"tdp-request-transfer",
}

type benchmarkLogsAction struct {
	BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN string
	RECEIVED_FORM_BACKEND                 string
	PUBLISH_TO                            string
	CONSUMING_TRANSACTION_QUEUE           string
	FO_USER_XDR_SUBMITTING_TO_BLOCKCHAIN  string
	PUBLISH_TO_BACKLINK                   string
	CONSUMING_BACKLINK_QUEUE              string
	TDP_REQUEST_SUBMITTED                 string
}

var BenchmarkLogsAction = benchmarkLogsAction{
	"backlink-xdr-submitting-to-blockchain",
	"received-form-backend",
	"publish-to-",
	"consuming-transaction-queue",
	"fo-user-xdr-submitting-to-blockchain",
	"publish-to-backlink",
	"consuming-backlink-queue",
	"tdp-request-submitted",
}

type benchmarkLogStatus struct {
	SUCCESS  string
	OK       string
	SENDING  string
	ERROR    string
	COMPLETE string
	START    string
}

var BenchmarkLogsStatus = benchmarkLogStatus{
	"success",
	"ok",
	"sending",
	"error",
	"complete",
	"start",
}
