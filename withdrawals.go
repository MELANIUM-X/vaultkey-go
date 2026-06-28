package vaultkey

import (
	"context"
	"net/url"
)

// WithdrawalsService handles developer-facing treasury-safe withdrawals.
type WithdrawalsService struct {
	client *Client
}

// Create creates a withdrawal request.
//
// VaultKey selects the payout wallet, applies treasury policy, and submits the
// underlying transfer when the withdrawal is ready to process.
func (w *WithdrawalsService) Create(ctx context.Context, payload CreateWithdrawalPayload) (Withdrawal, *ErrorResponse, error) {
	var resp Withdrawal
	apiErr, err := w.client.post(ctx, "/v2/withdrawals", payload, &resp)
	return resp, apiErr, err
}

// Get retrieves a withdrawal by ID.
func (w *WithdrawalsService) Get(ctx context.Context, withdrawalID string) (Withdrawal, *ErrorResponse, error) {
	var resp Withdrawal
	apiErr, err := w.client.get(ctx, "/v2/withdrawals/"+withdrawalID, &resp)
	return resp, apiErr, err
}

// List returns recent withdrawals. Pass zero-value options to use server defaults.
func (w *WithdrawalsService) List(ctx context.Context, opts ListWithdrawalsOptions) (WithdrawalList, *ErrorResponse, error) {
	q := url.Values{}
	if opts.UserID != "" {
		q.Set("user_id", opts.UserID)
	}
	if opts.After != "" {
		q.Set("after", opts.After)
	}

	path := "/v2/withdrawals"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var resp WithdrawalList
	apiErr, err := w.client.get(ctx, path, &resp)
	return resp, apiErr, err
}

// Cancel cancels a withdrawal while it is still before signing or broadcast.
func (w *WithdrawalsService) Cancel(ctx context.Context, withdrawalID string) (Withdrawal, *ErrorResponse, error) {
	var resp Withdrawal
	apiErr, err := w.client.post(ctx, "/v2/withdrawals/"+withdrawalID+"/cancel", nil, &resp)
	return resp, apiErr, err
}
