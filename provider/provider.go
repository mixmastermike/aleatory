package provider

import (
	"github.com/mixmastermike/aleatory/client"
)

// Provider provides the data factory structure
type Provider interface {

	// Message buffer
	Broadcast([]byte)

	// Un/register connections to the factory data
	Register(*client.Connection) error
	Unregister() (*client.Connection, error)

	// Inform the provider to begin generation
	Generate(GenerateParams)
	// Inform the provider to stop generation
	Stop()
}

// GenerateParams is a set of parameters used to inform the data generation
type GenerateParams interface{}
