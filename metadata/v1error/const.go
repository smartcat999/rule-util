package v1error

//metadata
const (
	ServiceIniting = iota + 2000
	ServiceStart
	ServiceStartFailed
	ServiceStop
	ServiceStopFailed
	SyncStart
	SyncStartFailed
	SyncStop
	SyncStopFailed
	Stream
	StreamError
	StreamOpen
	StreamOpenFail
	StreamClose
	StreamCloseFail
)
const (
	PubSubOK = iota + 2020
	PubSubError
	SubscriberOnline
	SubscriberOnlineFailed
	SubscriberOffline
	SubscriberOfflineFailed
	SubscriptionCreate
	SubscriptionCreateFailed
	SubscriptionDelete
	SubscriptionDeleteFailed
	SubscriptionOpen
	SubscriptionOpenFailed
	SubscriptionClose
	SubscriptionCloseFailed
)

const (
	RuleOK = iota + 2040
	RuleError
	RuleCreate
	RuleCreateFail
	RuleDelete
	RuleDeleteFail
	RuleOpen
	RuleOpenFail
	RuleClose
	RuleCloseFail
)

const (
	RuleAction = iota + 2060
	RuleActionError
	RuleActionCreate
	RuleActionCreateFail
	RuleActionDelete
	RuleActionDeleteFail
	RuleActionOpen
	RuleActionOpenFail
	RuleActionClose
	RuleActionCloseFail
)
