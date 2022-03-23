package kafka

// SourceOption is the function signature required to be considered an nats.SourceOption.
//type SourceOptionFn func(transport *sourcePostTransport) error

type SinkOptionFn func(transport *sinkTransport) error
