package client

// Option is the function signature required to be considered an client.Option.
type Option func(*ceClient) error
