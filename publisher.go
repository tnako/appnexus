package appnexus

import (
	"fmt"
	"net/http"
)

// PublisherService handles all requests to the publisher service API
type PublisherService struct {
	*Response
	client *Client
}

// Publisher is an audience publisher within the AppNexus console
type Publisher struct {
	ID                int    `json:"id,omitempty"`
	Code              string `json:"code,omitempty"`
	State             string `json:"state,omitempty"`
	Name              string `json:"name"`
	IsOO              bool   `json:"is_oo,omitempty"`
	ResellingExposure string `json:"reselling_exposure,omitempty"`
	BasePaymentRuleID int    `json:"base_payment_rule_id,omitempty"`
}

type publisherResponse struct {
	*http.Response
	Obj struct {
		Publisher  `json:"publisher,omitempty"`
		Publishers []Publisher `json:"publishers,omitempty"`
		Error      string      `json:"error"`
		Status     string      `json:"status"`
		Service    string      `json:"service"`
		Rate       Rate        `json:"dbg_info"`
	} `json:"response"`
}

// Get a publisher from the publisher service by ID
func (s *PublisherService) Get(publisherID int) (*Publisher, error) {

	path := fmt.Sprintf("publisher?id=%d", publisherID)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	r := &publisherResponse{}
	_, err = s.client.do(req, r)
	if err != nil {
		return nil, err
	}

	publisher := &r.Obj.Publisher
	return publisher, nil
}

// List available publishers from your AppNexus console
func (s *PublisherService) List(opt *ListOptions) ([]Publisher, *Response, error) {
	u, err := addOptions("publisher", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	publishers := &publisherResponse{}
	resp, err := s.client.do(req, publishers)
	if err != nil {
		return nil, resp, err
	}

	return publishers.Obj.Publishers, resp, err
}

// Add a new publisher
func (s *PublisherService) Add(item *Publisher) (*Response, error) {

	data := struct {
		Publisher `json:"publisher"`
	}{*item}

	req, err := s.client.newRequest("POST", "publisher?create_default_placement=false", data)

	if err != nil {
		return nil, err
	}

	result := &Response{}
	resp, err := s.client.do(req, result)
	if err != nil {
		return resp, err
	}

	item.ID = result.Obj.ID
	return result, nil
}

// Update an existing publisher with new data
// func (s *PublisherService) Update(memberID int, item Publisher) (*Response, error) {
//
// 	data := struct {
// 		Publisher `json:"publisher"`
// 	}{item}
//
// 	if item.ID < 1 {
// 		return nil, errors.New("Update Publisher requires a publisher to have an ID already")
// 	}
//
// 	req, err := s.client.newRequest("PUT", fmt.Sprintf("publisher/%d?id=%d", memberID, item.ID), data)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	result := &Response{}
// 	resp, err := s.client.do(req, result)
// 	if err != nil {
// 		return resp, err
// 	}
//
// 	return result, nil
// }

// Delete the specified publisher
// func (s *PublisherService) Delete(memberID int, item Publisher) error {
//
// 	data := struct {
// 		Publisher `json:"publisher"`
// 	}{item}
//
// 	if item.ID < 1 {
// 		return errors.New("Delete Publisher requires a publisher to have an ID already")
// 	}
//
// 	req, err := s.client.newRequest("DELETE", fmt.Sprintf("publisher/%d", memberID), data)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = s.client.do(req, nil)
// 	return err
// }