package appnexus

import (
	"errors"
	"fmt"
	"net/http"
)

// SiteService handles all requests to the site service API
type SiteService struct {
	*Response
	client *Client
}

// Site is an audience site within the AppNexus console
type Site struct {
	ID          int64  `json:"id,omitempty"`
	PublisherID int64  `json:"publisher_id"`
	Code        string `json:"code,omitempty"`
	State       string `json:"state,omitempty"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	SupplyType  string `json:"supply_type"`
}

type siteResponse struct {
	*http.Response
	Obj struct {
		Site    `json:"site,omitempty"`
		Sites   []Site `json:"sites,omitempty"`
		Error   string `json:"error"`
		Status  string `json:"status"`
		Service string `json:"service"`
		Rate    Rate   `json:"dbg_info"`
	} `json:"response"`
}

// Get a site from the site service by ID
func (s *SiteService) Get(params ...int64) (*Site, error) {
	var path string
	if len(params) > 1 {
		path = fmt.Sprintf("site?id=%d&publisher_id=%d", params[0], params[1])
	} else {
		path = fmt.Sprintf("site?id=%d", params[0])
	}
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	r := &siteResponse{}
	_, err = s.client.do(req, r)
	if err != nil {
		return nil, err
	}

	site := &r.Obj.Site
	return site, nil
}

// List available sites from your AppNexus console
func (s *SiteService) List() ([]Site, *Response, error) {
	req, err := s.client.newRequest("GET", "site", nil)
	if err != nil {
		return nil, nil, err
	}

	sites := &siteResponse{}
	resp, err := s.client.do(req, sites)
	if err != nil {
		return nil, resp, err
	}

	return sites.Obj.Sites, resp, err
}

// Add a new site
func (s *SiteService) Add(item *Site, pubID int64) (*Response, error) {

	data := struct {
		Site `json:"site"`
	}{*item}

	req, err := s.client.newRequest("POST", fmt.Sprintf("site?publisher_id=%d", pubID), data)

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

// Update an existing site with new data
func (s *SiteService) Update(item Site, pubID int64) (*Response, error) {

	data := struct {
		Site `json:"site"`
	}{item}

	if item.ID < 1 {
		return nil, errors.New("Update Site requires a site to have an ID already")
	}

	req, err := s.client.newRequest("PUT", fmt.Sprintf("site?id=%d&publisher_id=%d", item.ID, pubID), data)

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

// Delete the specified site
func (s *SiteService) Delete(siteID int64, pubID int64) error {
	req, err := s.client.newRequest("DELETE", fmt.Sprintf("site?id=%d&publisher_id=%d", siteID, pubID), nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
