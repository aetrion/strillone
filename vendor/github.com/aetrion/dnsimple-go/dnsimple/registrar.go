package dnsimple

import (
	"fmt"
)

// RegistrarService handles communication with the registrar related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/registrar/
type RegistrarService struct {
	client *Client
}

// DomainCheck represents the result of a domain check.
type DomainCheck struct {
	Domain    string `json:"domain"`
	Available bool   `json:"available"`
	Premium   bool   `json:"premium"`
}

// DomainCheckResponse represents a response from the domain check.
type DomainCheckResponse struct {
	Response
	Data *DomainCheck `json:"data"`
}

// CheckDomain checks a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#check
func (s *RegistrarService) CheckDomain(accountID string, domainName string) (*DomainCheckResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/check", accountID, domainName))
	checkResponse := &DomainCheckResponse{}

	resp, err := s.client.get(path, checkResponse)
	if err != nil {
		return nil, err
	}

	checkResponse.HttpResponse = resp
	return checkResponse, nil
}

// DomainRegisterRequest represents the attributes you can pass to a register API request.
// Some attributes are mandatory.
type DomainRegisterRequest struct {
	// The ID of the Contact to use as registrant for the domain
	RegistrantID int `json:"registrant_id"`
	// Set to true to enable the whois privacy service. An extra cost may apply.
	// Default to false.
	EnableWhoisPrivacy bool `json:"private_whois,omitempty"`
	// Set to true to enable the auto-renewal of the domain.
	// Default to true.
	EnableAutoRenewal bool `json:"auto_renew,omitempty"`
}

// DomainRegistrationResponse represents a response from an API method that results in a domain registration.
type DomainRegistrationResponse struct {
	Response
	Data *Domain `json:"data"`
}

// RegisterDomain registers a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#register
func (s *RegistrarService) RegisterDomain(accountID string, domainName string, request *DomainRegisterRequest) (*DomainRegistrationResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/registration", accountID, domainName))
	registrationResponse := &DomainRegistrationResponse{}

	// TODO: validate mandatory attributes RegistrantID

	resp, err := s.client.post(path, request, registrationResponse)
	if err != nil {
		return nil, err
	}

	registrationResponse.HttpResponse = resp
	return registrationResponse, nil
}

// DomainTransferRequest represents the attributes you can pass to a transfer API request.
// Some attributes are mandatory.
type DomainTransferRequest struct {
	// The ID of the Contact to use as registrant for the domain
	RegistrantID int `json:"registrant_id"`
	// The Auth-Code required to transfer the domain.
	// This is provided by the current registrar of the domain.
	AuthInfo string `json:"auth_info,omitempty"`
	// Set to true to enable the whois privacy service. An extra cost may apply.
	// Default to false.
	EnableWhoisPrivacy bool `json:"private_whois,omitempty"`
	// Set to true to enable the auto-renewal of the domain.
	// Default to true.
	EnableAutoRenewal bool `json:"auto_renew,omitempty"`
}

// DomainTransferResponse represents a response from an API method that results in a domain transfer.
type DomainTransferResponse struct {
	Response
	Data *Domain `json:"data"`
}

// TransferDomain transfers a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#transfer
func (s *RegistrarService) TransferDomain(accountID string, domainName string, request *DomainTransferRequest) (*DomainTransferResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/transfer", accountID, domainName))
	transferResponse := &DomainTransferResponse{}

	// TODO: validate mandatory attributes RegistrantID

	resp, err := s.client.post(path, request, transferResponse)
	if err != nil {
		return nil, err
	}

	transferResponse.HttpResponse = resp
	return transferResponse, nil
}

// DomainTransferOutResponse represents a response from an API method that results in a domain transfer out.
type DomainTransferOutResponse struct {
	Response
	Data *Domain `json:"data"`
}

// Transfer out a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#transfer-out
func (s *RegistrarService) TransferDomainOut(accountID string, domainName string) (*DomainTransferOutResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/transfer_out", accountID, domainName))
	transferResponse := &DomainTransferOutResponse{}

	resp, err := s.client.post(path, nil, nil)
	if err != nil {
		return nil, err
	}

	transferResponse.HttpResponse = resp
	return transferResponse, nil
}

// DomainRenewRequest represents the attributes you can pass to a renew API request.
// Some attributes are mandatory.
type DomainRenewRequest struct {
	// The number of years
	Period int `json:"period"`
}

// DomainRenewalResponse represents a response from an API method that results in a domain renewal.
type DomainRenewalResponse struct {
	Response
	Data *Domain `json:"data"`
}

// RenewDomain renews a domain name.
//
// See https://developer.dnsimple.com/v2/registrar/#register
func (s *RegistrarService) RenewDomain(accountID string, domainName string, request *DomainRenewRequest) (*DomainRenewalResponse, error) {
	path := versioned(fmt.Sprintf("/%v/registrar/domains/%v/renewal", accountID, domainName))
	renewalResponse := &DomainRenewalResponse{}

	resp, err := s.client.post(path, request, renewalResponse)
	if err != nil {
		return nil, err
	}

	renewalResponse.HttpResponse = resp
	return renewalResponse, nil
}
