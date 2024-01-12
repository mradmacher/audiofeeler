package audiofeeler

import (
	"gopkg.in/yaml.v3"
	"io"
)

func LoadEvents(reader io.Reader) (events []Event, err error) {
	decoder := yaml.NewDecoder(reader)
	err = decoder.Decode(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}
