package vaultkey

import "context"

// MonitorService handles wallet transaction monitoring subscriptions.
type MonitorService struct {
	client *Client
}

// Subscribe opts a wallet into transaction monitoring.
func (m *MonitorService) Subscribe(ctx context.Context, payload MonitorSubscriptionPayload) (MonitorSubscription, *ErrorResponse, error) {
	var resp MonitorSubscription
	apiErr, err := m.client.post(ctx, "/monitor/wallets", payload, &resp)
	return resp, apiErr, err
}

// List returns active monitoring subscriptions for the project.
func (m *MonitorService) List(ctx context.Context) ([]MonitorSubscription, *ErrorResponse, error) {
	var resp []MonitorSubscription
	apiErr, err := m.client.get(ctx, "/monitor/wallets", &resp)
	return resp, apiErr, err
}

// Get retrieves a monitoring subscription by ID.
func (m *MonitorService) Get(ctx context.Context, subscriptionID string) (MonitorSubscription, *ErrorResponse, error) {
	var resp MonitorSubscription
	apiErr, err := m.client.get(ctx, "/monitor/wallets/"+subscriptionID, &resp)
	return resp, apiErr, err
}

// Unsubscribe removes a wallet monitoring subscription.
func (m *MonitorService) Unsubscribe(ctx context.Context, subscriptionID string) (MonitorUnsubscribeResult, *ErrorResponse, error) {
	var resp MonitorUnsubscribeResult
	apiErr, err := m.client.delete(ctx, "/monitor/wallets/"+subscriptionID, nil, &resp)
	return resp, apiErr, err
}
