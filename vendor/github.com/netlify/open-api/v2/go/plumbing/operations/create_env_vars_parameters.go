// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/netlify/open-api/v2/go/models"
)

// NewCreateEnvVarsParams creates a new CreateEnvVarsParams object
// with the default values initialized.
func NewCreateEnvVarsParams() *CreateEnvVarsParams {
	var ()
	return &CreateEnvVarsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateEnvVarsParamsWithTimeout creates a new CreateEnvVarsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateEnvVarsParamsWithTimeout(timeout time.Duration) *CreateEnvVarsParams {
	var ()
	return &CreateEnvVarsParams{

		timeout: timeout,
	}
}

// NewCreateEnvVarsParamsWithContext creates a new CreateEnvVarsParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateEnvVarsParamsWithContext(ctx context.Context) *CreateEnvVarsParams {
	var ()
	return &CreateEnvVarsParams{

		Context: ctx,
	}
}

// NewCreateEnvVarsParamsWithHTTPClient creates a new CreateEnvVarsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateEnvVarsParamsWithHTTPClient(client *http.Client) *CreateEnvVarsParams {
	var ()
	return &CreateEnvVarsParams{
		HTTPClient: client,
	}
}

/*CreateEnvVarsParams contains all the parameters to send to the API endpoint
for the create env vars operation typically these are written to a http.Request
*/
type CreateEnvVarsParams struct {

	/*AccountID
	  Scope response to account_id

	*/
	AccountID string
	/*EnvVars*/
	EnvVars []*models.CreateEnvVarsParamsBodyItems
	/*SiteID
	  If provided, create an environment variable on the site level, not the account level

	*/
	SiteID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create env vars params
func (o *CreateEnvVarsParams) WithTimeout(timeout time.Duration) *CreateEnvVarsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create env vars params
func (o *CreateEnvVarsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create env vars params
func (o *CreateEnvVarsParams) WithContext(ctx context.Context) *CreateEnvVarsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create env vars params
func (o *CreateEnvVarsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create env vars params
func (o *CreateEnvVarsParams) WithHTTPClient(client *http.Client) *CreateEnvVarsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create env vars params
func (o *CreateEnvVarsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAccountID adds the accountID to the create env vars params
func (o *CreateEnvVarsParams) WithAccountID(accountID string) *CreateEnvVarsParams {
	o.SetAccountID(accountID)
	return o
}

// SetAccountID adds the accountId to the create env vars params
func (o *CreateEnvVarsParams) SetAccountID(accountID string) {
	o.AccountID = accountID
}

// WithEnvVars adds the envVars to the create env vars params
func (o *CreateEnvVarsParams) WithEnvVars(envVars []*models.CreateEnvVarsParamsBodyItems) *CreateEnvVarsParams {
	o.SetEnvVars(envVars)
	return o
}

// SetEnvVars adds the envVars to the create env vars params
func (o *CreateEnvVarsParams) SetEnvVars(envVars []*models.CreateEnvVarsParamsBodyItems) {
	o.EnvVars = envVars
}

// WithSiteID adds the siteID to the create env vars params
func (o *CreateEnvVarsParams) WithSiteID(siteID *string) *CreateEnvVarsParams {
	o.SetSiteID(siteID)
	return o
}

// SetSiteID adds the siteId to the create env vars params
func (o *CreateEnvVarsParams) SetSiteID(siteID *string) {
	o.SiteID = siteID
}

// WriteToRequest writes these params to a swagger request
func (o *CreateEnvVarsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param account_id
	if err := r.SetPathParam("account_id", o.AccountID); err != nil {
		return err
	}

	if o.EnvVars != nil {
		if err := r.SetBodyParam(o.EnvVars); err != nil {
			return err
		}
	}

	if o.SiteID != nil {

		// query param site_id
		var qrSiteID string
		if o.SiteID != nil {
			qrSiteID = *o.SiteID
		}
		qSiteID := qrSiteID
		if qSiteID != "" {
			if err := r.SetQueryParam("site_id", qSiteID); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
