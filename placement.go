package appnexus

import (
	"errors"
	"fmt"
	"net/http"
)

// PlacementService handles all requests to the placement service API
type PlacementService struct {
	*Response
	client *Client
}

// Placement is an audience placement within the AppNexus console
type Placement struct {
	ID          int64  `json:"id,omitempty"`
	PublisherID int64  `json:"publisher_id"`
	SiteID      int64  `json:"site_id,omitempty"`
	Code        string `json:"code"`
	State       string `json:"state,omitempty"`
	Name        string `json:"name"`
}

type placementResponse struct {
	*http.Response
	Obj struct {
		Placement  `json:"placement,omitempty"`
		Placements []Placement `json:"placements,omitempty"`
		Error      string      `json:"error"`
		Status     string      `json:"status"`
		Service    string      `json:"service"`
		Rate       Rate        `json:"dbg_info"`
	} `json:"response"`
}

// Get a placement from the placement service by ID
func (s *PlacementService) Get(placementID int64) (*Placement, error) {
	path := fmt.Sprintf("placement?id=%d", placementID)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	r := &placementResponse{}
	_, err = s.client.do(req, r)
	if err != nil {
		return nil, err
	}

	placement := &r.Obj.Placement
	return placement, nil
}

// List available placements from your AppNexus console
func (s *PlacementService) List(pubID int64) ([]Placement, *Response, error) {
	path := fmt.Sprintf("placement?publisher_id=%d", pubID)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	placements := &placementResponse{}
	resp, err := s.client.do(req, placements)
	if err != nil {
		return nil, resp, err
	}

	return placements.Obj.Placements, resp, err
}

// Add a new placement
func (s *PlacementService) Add(item *Placement) (*Response, error) {

	data := struct {
		Placement `json:"placement"`
	}{*item}

	req, err := s.client.newRequest("POST", fmt.Sprintf("placement?publisher_id=%d", item.PublisherID), data)

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

// Update an existing placement with new data
func (s *PlacementService) Update(item Placement) (*Response, error) {

	data := struct {
		Placement `json:"placement"`
	}{item}

	if item.ID < 1 {
		return nil, errors.New("Update Placement requires a placement to have an ID already")
	}

	req, err := s.client.newRequest("PUT", fmt.Sprintf("placement?id=%d&publisher_id=%d", item.ID, item.PublisherID), data)

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

// Delete the specified placement
func (s *PlacementService) Delete(placementID int64, pubID int64) error {
	req, err := s.client.newRequest("DELETE", fmt.Sprintf("placement?id=%d&publisher_id=%d", placementID, pubID), nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
