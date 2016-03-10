package dnsimple

import (
	"fmt"
)

// ContactsService handles communication with the contact related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/contacts/
type ContactsService struct {
	client *Client
}

// Contact represents a Contact in DNSimple.
type Contact struct {
	ID            int    `json:"id,omitempty"`
	AccountID     int    `json:"account_id,omitempty"`
	Label         string `json:"label,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	JobTitle      string `json:"job_title,omitempty"`
	Organization  string `json:"organization_name,omitempty"`
	Address1      string `json:"address1,omitempty"`
	Address2      string `json:"address2,omitempty"`
	City          string `json:"city,omitempty"`
	StateProvince string `json:"state_province,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Country       string `json:"country,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Fax           string `json:"fax,omitempty"`
	Email         string `json:"email_address,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

func contactPath(accountID string, contact interface{}) string {
	if contact != nil {
		return fmt.Sprintf("/%v/contacts/%v", accountID, contact)
	}
	return fmt.Sprintf("/%v/contacts", accountID)
}

// ContactResponse represents a response from an API method that returns a Contact struct.
type ContactResponse struct {
	Response
	Data *Contact `json:"data"`
}

// ContactsResponse represents a response from an API method that returns a collection of Contact struct.
type ContactsResponse struct {
	Response
	Data []Contact `json:"data"`
}

// ListContacts list the contacts for an account.
//
// See https://developer.dnsimple.com/v2/contacts/#list
func (s *ContactsService) ListContacts(accountID string, options *ListOptions) (*ContactsResponse, error) {
	path := versioned(contactPath(accountID, nil))
	contactsResponse := &ContactsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, contactsResponse)
	if err != nil {
		return contactsResponse, err
	}

	contactsResponse.HttpResponse = resp
	return contactsResponse, nil
}

// CreateContact creates a new contact.
//
// See https://developer.dnsimple.com/v2/contacts/#create
func (s *ContactsService) CreateContact(accountID string, contactAttributes Contact) (*ContactResponse, error) {
	path := versioned(contactPath(accountID, nil))
	contactResponse := &ContactResponse{}

	resp, err := s.client.post(path, contactAttributes, contactResponse)
	if err != nil {
		return nil, err
	}

	contactResponse.HttpResponse = resp
	return contactResponse, nil
}

// GetContact fetches a contact.
//
// See https://developer.dnsimple.com/v2/contacts/#get
func (s *ContactsService) GetContact(accountID string, contactID int) (*ContactResponse, error) {
	path := versioned(contactPath(accountID, contactID))
	contactResponse := &ContactResponse{}

	resp, err := s.client.get(path, contactResponse)
	if err != nil {
		return nil, err
	}

	contactResponse.HttpResponse = resp
	return contactResponse, nil
}

// UpdateContact updates a contact.
//
// See https://developer.dnsimple.com/v2/contacts/#update
func (s *ContactsService) UpdateContact(accountID string, contactID int, contactAttributes Contact) (*ContactResponse, error) {
	path := versioned(contactPath(accountID, contactID))
	contactResponse := &ContactResponse{}

	resp, err := s.client.patch(path, contactAttributes, contactResponse)
	if err != nil {
		return nil, err
	}

	contactResponse.HttpResponse = resp
	return contactResponse, nil
}

// DeleteContact PERMANENTLY deletes a contact from the account.
//
// See https://developer.dnsimple.com/v2/contacts/#delete
func (s *ContactsService) DeleteContact(accountID string, contactID int) (*ContactResponse, error) {
	path := versioned(contactPath(accountID, contactID))
	contactResponse := &ContactResponse{}

	resp, err := s.client.delete(path, nil, nil)
	if err != nil {
		return nil, err
	}

	contactResponse.HttpResponse = resp
	return contactResponse, nil
}
