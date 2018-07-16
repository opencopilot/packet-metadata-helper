package packetmd

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

// AddTag adds a tag to a packet device
func AddTag(client *packngo.Client, deviceID, tag string) error {
	device, _, err := client.Devices.Get(deviceID)
	if err != nil {
		return err
	}
	tags := append(device.Tags, tag)
	client.Devices.Update(deviceID, &packngo.DeviceUpdateRequest{
		Tags: &tags,
	})

	return nil
}

// RemoveTag removes a tag from a packet device
func RemoveTag(client *packngo.Client, deviceID, tag string) error {
	device, _, err := client.Devices.Get(deviceID)
	if err != nil {
		return err
	}
	tags := make([]string, 0)
	for _, t := range device.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}
	client.Devices.Update(deviceID, &packngo.DeviceUpdateRequest{
		Tags: &tags,
	})
	return nil
}

// UpdateTag updates an existing tag on a packet device. If it can't find the tag, nothing happens
func UpdateTag(client *packngo.Client, deviceID, tag, newTag string) error {
	device, _, err := client.Devices.Get(deviceID)
	if err != nil {
		return err
	}
	tags := make([]string, 0)
	for _, t := range device.Tags {
		if t == tag {
			t = newTag
		}
		tags = append(tags, t)
	}
	client.Devices.Update(deviceID, &packngo.DeviceUpdateRequest{
		Tags: &tags,
	})
	return nil
}
