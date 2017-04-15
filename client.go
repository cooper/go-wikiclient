// Copyright (c) 2017, Mitchell Cooper
package wikiclient

import (
	"errors"
	"strconv"
	"time"
)

// a Client is formed by pairing a transport with a session
type Client struct {
	Transport Transport     // wikiclient transport
	Session   *Session      // wikiclient session
	Timeout   time.Duration // how long to waits on requests
}

// create a client and clean the session if necessary
func NewClient(tr Transport, sess *Session, timeout time.Duration) (c Client) {
	c = Client{tr, sess, timeout}
	c.Clean()
	return
}

// display a page
func (c Client) DisplayPage(pageName string) (Message, error) {
	return c.Request("page", map[string]interface{}{"name": pageName})
}

// display an image
func (c Client) DisplayImage(imageName string, width, height int) (Message, error) {
	return c.Request("image", map[string]interface{}{
		"name":   imageName,
		"width":  strconv.Itoa(width),
		"height": strconv.Itoa(height),
	})
}

// display category posts
func (c Client) DisplayCategoryPosts(categoryName string, pageN int) (Message, error) {
	if pageN <= 0 {
		pageN = 1
	}
	return c.Request("cat_posts", map[string]interface{}{
		"name":   categoryName,
		"page_n": string(pageN),
	})
}

// send a message and block until we get its response
func (c Client) Request(command string, args messageArgs) (Message, error) {
	return c.RequestMessage(NewMessage(command, args))
}

// send a message and block until we get its response
func (c Client) RequestMessage(req Message) (Message, error) {
	return c.requestGeneric(false, req, false)
}

// same as Request() except that it does not Connect()
func (c Client) requestConnecting(command string, args messageArgs) (Message, error) {
	return c.requestGeneric(true, NewMessage(command, args), true)
}

func (c Client) requestGeneric(connecting bool, req Message, neverTimeout bool) (res Message, err error) {
	if !connecting {
		c.Connect()
	}

	// send
	if err = c.sendMessage(req); err != nil {
		return
	}
	
	// timeout channel
	var timeout <-chan time.Time
	if !neverTimeout {
		timeout = time.After(c.Timeout)
	}

	// await the response, or give up after the timeout
	select {
	case res = <-c.Transport.readMessages():

		// this is the correct ID
		if res.ID == req.ID {
			return
		}

		// some other message
		err = errors.New("Got response with incorrect message ID")
		return

	case <-timeout:
		err = errors.New("Timed out")
		return
	}
}

// send a message to the transport, but do not await reply
func (c Client) sendMessage(msg Message) error {

	// the transport is dead!
	if c.Transport.Dead() {
		return errors.New("transport is dead")
	}

	return c.Transport.writeMessage(msg)
}

// authenticate for read/write access as necessary
func (c Client) Connect() error {
	c.Clean()

	// the transport is not authenticated
	if !c.Session.ReadAccess {
		wikires, err := c.requestConnecting("wiki", map[string]interface{}{
			"name":     c.Session.WikiName,
			"password": c.Session.WikiPassword,
			"config":   true,
		})
		if err != nil {
			return err
		}
		info, ok := wikires.Args["config"].(map[string]interface{})
		if !ok {
			return errors.New("wikiserver did not provide wiki configuration")
		}
		c.Session.Config = info
		c.Session.ReadAccess = true
		c.Transport.SelectWiki(c.Session.WikiName)
	}

	// TODO: if the transport is not write authenticated and we have
	// credentials in the session, send them now

	// this wiki is not selected, so send select
	if c.Transport.SelectWiki(c.Session.WikiName) {
		_, err := c.requestConnecting("select", map[string]interface{}{
			"name": c.Session.WikiName,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Client) Clean() {
	c.Session.clean(c.Transport)
}
