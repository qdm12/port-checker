package config

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

type Settings struct {
	ListeningAddress *string
	RootURL          *string
}

func (s *Settings) SetDefaults() {
	s.ListeningAddress = gosettings.DefaultPointer(s.ListeningAddress, ":8000")
	s.RootURL = gosettings.DefaultPointer(s.RootURL, "/")
}

func (s Settings) Validate() (err error) {
	return validate.ListeningAddress(*s.ListeningAddress, os.Getuid())
}

func (s Settings) String() string {
	return s.toLinesNode().String()
}

func (s Settings) toLinesNode() *gotree.Node {
	node := gotree.New("Settings")
	node.Appendf("Listening address: %s", *s.ListeningAddress)
	node.Appendf("Root URL: %s", *s.RootURL)
	return node
}

func ptrTo[T any](value T) *T { return &value }

func (s *Settings) Read(reader *reader.Reader) (err error) {
	// Retro-compatibility
	port, err := reader.Uint16Ptr("LISTENING_PORT")
	if err != nil {
		return err
	}
	if port != nil { // Retro-compatibility
		s.ListeningAddress = ptrTo(fmt.Sprintf(":%d", *port))
	} else {
		s.ListeningAddress = reader.Get("LISTENING_ADDRESS")
	}

	s.RootURL = reader.Get("ROOT_URL")

	return nil
}
