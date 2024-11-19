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
	ListeningPort *uint16
	RootURL       *string
}

func (s *Settings) SetDefaults() {
	const defaultListeningPort = 8000
	s.ListeningPort = gosettings.DefaultPointer(s.ListeningPort, defaultListeningPort)
	s.RootURL = gosettings.DefaultPointer(s.RootURL, "/")
}

func (s Settings) Validate() (err error) {
	address := fmt.Sprintf(":%d", *s.ListeningPort)
	return validate.ListeningAddress(address, os.Getuid())
}

func (s Settings) String() string {
	return s.toLinesNode().String()
}

func (s Settings) toLinesNode() *gotree.Node {
	node := gotree.New("Settings")
	node.Appendf("Listening port: %d", *s.ListeningPort)
	node.Appendf("Root URL: %s", *s.RootURL)
	return node
}

func (s *Settings) Read(reader *reader.Reader) (err error) {
	s.ListeningPort, err = reader.Uint16Ptr("LISTENING_PORT")
	if err != nil {
		return err
	}

	s.RootURL = reader.Get("ROOT_URL")

	return nil
}
