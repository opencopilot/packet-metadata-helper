package packetmetadatahelper

import (
	"regexp"
	"strings"

	"github.com/packethost/packngo"
)

// KV represents a key-value pair stored in a device tag
type KV struct {
	Key   string
	Value string
}

// GetKVPairs will parse a device's tags looking for tags of the form `key=value` where the delimiter (=) can be set to an arbitrary string
func GetKVPairs(device *packngo.Device, delimiter string) ([]*KV, error) {
	tags := device.Tags
	pairs := make([]*KV, 0)

	if delimiter == "" {
		delimiter = "="
	}

	for _, tag := range tags { // loop through the device tags
		// match any tag that is of the form key=value
		matched, err := regexp.MatchString(".+"+delimiter+".+", tag)
		if err != nil {
			return nil, err
		}
		// if the tag doesn't match this, skip it
		if !matched {
			continue
		}
		// split the tag on the delimiter, but only into two parts
		split := strings.SplitN(tag, delimiter, 2)
		if len(split) < 2 {
			continue
		}
		pairs = append(pairs, &KV{
			Key:   split[0],
			Value: split[1],
		})
	}

	return pairs, nil
}
