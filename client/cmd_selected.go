package client

import (
	"errors"

	imap "github.com/emersion/imap/common"
	"github.com/emersion/imap/commands"
	"github.com/emersion/imap/responses"
)

// TODO: CHECK

// Permanently removes all messages that have the \Deleted flag set from the
// currently selected mailbox, and returns to the authenticated state from the
// selected state.
func (c *Client) Close() (err error) {
	if c.State != imap.SelectedState {
		err = errors.New("No mailbox selected")
		return
	}

	cmd := &commands.Close{}

	status, err := c.execute(cmd, nil)
	if err != nil {
		return
	}

	err = status.Err()
	if err != nil {
		return
	}

	c.State = imap.AuthenticatedState
	c.Selected = nil
	return
}

// Permanently removes all messages that have the \Deleted flag set from the
// currently selected mailbox.
// If ch is not nil, sends sequence IDs of each deleted message to this channel.
func (c *Client) Expunge(ch chan<- uint32) (err error) {
	defer close(ch)

	if c.State != imap.SelectedState {
		err = errors.New("No mailbox selected")
		return
	}

	cmd := &commands.Expunge{}

	var res *responses.Expunge
	if ch != nil {
		res = &responses.Expunge{SeqIds: ch}
	}

	status, err := c.execute(cmd, res)
	if err != nil {
		return
	}

	err = status.Err()
	return
}

// Searches the mailbox for messages that match the given searching criteria.
// Searching criteria consist of one or more search keys. The response contains
// a list of message sequence IDs corresponding to those messages that match the
// searching criteria.
// When multiple keys are specified, the result is the intersection (AND
// function) of all the messages that match those keys.
// Criteria must be UTF-8 encoded.
// See RFC 3501 section 6.4.4 for a list of searching criteria.
func (c *Client) Search(criteria []interface{}) (ids []uint32, err error) {
	if c.State != imap.SelectedState {
		err = errors.New("No mailbox selected")
		return
	}

	cmd := &commands.Search{
		Charset: "UTF-8",
		Criteria: criteria,
	}

	res := &responses.Search{}

	status, err := c.execute(cmd, res)
	if err != nil {
		return
	}

	err = status.Err()
	ids = res.Ids
	return
}

// Retrieves data associated with a message in the mailbox.
// See RFC 3501 section 6.4.5 for a list of items that can be requested.
func (c *Client) Fetch(seqset *imap.SeqSet, items []string, ch chan<- *imap.Message) (err error) {
	defer close(ch)

	if c.State != imap.SelectedState {
		err = errors.New("No mailbox selected")
		return
	}

	cmd := &commands.Fetch{
		SeqSet: seqset,
		Items: items,
	}
	res := &responses.Fetch{Messages: ch}

	status, err := c.execute(cmd, res)
	if err != nil {
		return
	}

	err = status.Err()
	return
}

// TODO: STORE, COPY, UID