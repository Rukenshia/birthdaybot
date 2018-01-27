package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// SendRocketchatMessage sends a Message to a Rocket.Chat server
func SendRocketchatMessage(url, channel, message string) error {
	var c http.Client

	data, err := json.Marshal(map[string]interface{}{
		"text":    message,
		"channel": channel,
	})

	buf := bytes.NewBuffer(data)

	res, err := c.Post(url, "application/json", buf)
	if res.StatusCode != 200 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("Remote response with status %d: %s", res.StatusCode, string(data)))
	}
	log.Debugf("%v", res)
	return err
}
