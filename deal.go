package appnexus

import (
	"errors"
	"fmt"
	"net/http"
)

// DealService handles all requests to the deal service API
type DealService struct {
	*Response
	client *Client
}

// Type is a nested part of deal within the AppNexus console
type Type struct {
	ID int64 `json:"id,omitempty"`
}

// AuctionType is a nested part of deal within the AppNexus console
type AuctionType struct {
	ID int64 `json:"id,omitempty"`
}

// Buyer is a nested part of deal within the AppNexus console
type Buyer struct {
	ID int64 `json:"id,omitempty"`
}

// Deal is an audience deal within the AppNexus console
type Deal struct {
	ID          int64        `json:"id,omitempty"`
	FloorPrice  float64      `json:"floor_price,omitempty"`
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Active      bool         `json:"active"`
	StartDate   string       `json:"start_date,omitempty"`
	EndDate     string       `json:"end_date,omitempty"`
	Type        *Type        `json:"type,omitempty"`
	AuctionType *AuctionType `json:"auction_type,omitempty"`
	Buyer       *Buyer       `json:"buyer,omitempty"`
}

type dealResponse struct {
	*http.Response
	Obj struct {
		Deal    `json:"deal,omitempty"`
		Deals   []Deal `json:"deals,omitempty"`
		Error   string `json:"error"`
		Status  string `json:"status"`
		Service string `json:"service"`
		Rate    Rate   `json:"dbg_info"`
	} `json:"response"`
}

// Get a deal from the deal service by ID
func (s *DealService) Get(dealID int64) (*Deal, error) {
	path := fmt.Sprintf("deal?id=%d", dealID)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	r := &dealResponse{}
	_, err = s.client.do(req, r)
	if err != nil {
		return nil, err
	}

	deal := &r.Obj.Deal
	return deal, nil
}

// List available deals from your AppNexus console
func (s *DealService) List() ([]Deal, *Response, error) {
	req, err := s.client.newRequest("GET", "deal", nil)
	if err != nil {
		return nil, nil, err
	}

	deals := &dealResponse{}
	resp, err := s.client.do(req, deals)
	if err != nil {
		return nil, resp, err
	}

	return deals.Obj.Deals, resp, err
}

// Add a new deal
func (s *DealService) Add(item *Deal) (*Response, error) {

	data := struct {
		Deal `json:"deal"`
	}{*item}

	req, err := s.client.newRequest("POST", "deal", data)

	if err != nil {
		return nil, err
	}

	result := &Response{}
	resp, err := s.client.do(req, result)
	if err != nil {
		return resp, err
	}

	item.ID, _ = result.Obj.ID.Int64()
	return result, nil
}

// Update an existing deal with new data
func (s *DealService) Update(item Deal) (*Response, error) {

	data := struct {
		Deal `json:"deal"`
	}{item}

	if item.ID < 1 {
		return nil, errors.New("Update Deal requires a deal to have an ID already")
	}

	req, err := s.client.newRequest("PUT", fmt.Sprintf("deal?id=%d", item.ID), data)

	if err != nil {
		return nil, err
	}

	result := &Response{}
	resp, err := s.client.do(req, result)
	if err != nil {
		return resp, err
	}

	return result, nil
}

// Delete the specified deal
func (s *DealService) Delete(dealID int64) error {
	req, err := s.client.newRequest("DELETE", fmt.Sprintf("deal?id=%d", dealID), nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
