package vaultkey

// ErrorResponse is returned by the API on non-2xx responses.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return e.Code + ": " + e.Message
}

// ChainType is either "evm" or "solana".
type ChainType string

const (
	ChainTypeEVM    ChainType = "evm"
	ChainTypeSolana ChainType = "solana"
	ChainTypeTron   ChainType = "tron"
)

// TransferSpeed controls transaction priority for stablecoin transfers.
type TransferSpeed string

const (
	SpeedSlow   TransferSpeed = "slow"
	SpeedNormal TransferSpeed = "normal"
	SpeedFast   TransferSpeed = "fast"
)

// JobStatus is the current state of an async job.
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

// WithdrawalStatus is the lifecycle state for a withdrawal request.
type WithdrawalStatus string

const (
	WithdrawalStatusCreated         WithdrawalStatus = "created"
	WithdrawalStatusPendingApproval WithdrawalStatus = "pending_approval"
	WithdrawalStatusApproved        WithdrawalStatus = "approved"
	WithdrawalStatusFundingRequired WithdrawalStatus = "funding_required"
	WithdrawalStatusFunding         WithdrawalStatus = "funding"
	WithdrawalStatusSigning         WithdrawalStatus = "signing"
	WithdrawalStatusCompleted       WithdrawalStatus = "completed"
	WithdrawalStatusFailed          WithdrawalStatus = "failed"
	WithdrawalStatusCancelled       WithdrawalStatus = "cancelled"
	WithdrawalStatusManualReview    WithdrawalStatus = "manual_review"
	WithdrawalStatusRejected        WithdrawalStatus = "rejected"
)

// MonitorDirection controls which transaction directions trigger webhooks.
type MonitorDirection string

const (
	MonitorDirectionIncoming MonitorDirection = "incoming"
	MonitorDirectionOutgoing MonitorDirection = "outgoing"
)

// MonitorSubscriptionStatus is the lifecycle state of a monitor subscription.
type MonitorSubscriptionStatus string

const (
	MonitorSubscriptionStatusActive   MonitorSubscriptionStatus = "active"
	MonitorSubscriptionStatusInactive MonitorSubscriptionStatus = "inactive"
)

// ── Chains ────────────────────────────────────────────────────────────────────

// Chain represents a supported EVM chain.
type Chain struct {
	Name         string `json:"name"`
	ChainID      string `json:"chain_id"`
	NativeSymbol string `json:"native_symbol"`
	LegacySymbol string `json:"legacy_symbol,omitempty"`
	Testnet      bool   `json:"testnet"`
}

// ── Wallets ───────────────────────────────────────────────────────────────────

// Wallet represents a custodial wallet.
type Wallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ChainType ChainType `json:"chain_type"`
	Address   string    `json:"address"`
	Label     string    `json:"label,omitempty"`
	CreatedAt string    `json:"created_at"`
}

// CreateWalletPayload is the request body for wallet creation.
type CreateWalletPayload struct {
	UserID    string    `json:"user_id"`
	ChainType ChainType `json:"chain_type"`
	Label     string    `json:"label,omitempty"`
}

// WalletList is the paginated response from listing wallets.
type WalletList struct {
	Wallets    []Wallet `json:"wallets"`
	NextCursor string   `json:"next_cursor,omitempty"`
	HasMore    bool     `json:"has_more"`
}

// ── Signing ───────────────────────────────────────────────────────────────────

// SigningJob is returned by async signing operations.
type SigningJob struct {
	JobID  string    `json:"job_id"`
	Status JobStatus `json:"status"`
}

// SignMessagePayload is the request body for signing a message.
type SignMessagePayload struct {
	Payload        map[string]any `json:"payload"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
}

// ── Balance ───────────────────────────────────────────────────────────────────

// EVMBalance is the response from an EVM balance lookup.
type EVMBalance struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	RawBalance string `json:"raw_balance"`
	Symbol     string `json:"symbol"`
	ChainName  string `json:"chain_name"`
	ChainID    string `json:"chain_id"`
}

// SolanaBalance is the response from a Solana balance lookup.
type SolanaBalance struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	RawBalance string `json:"raw_balance"`
	Symbol     string `json:"symbol"`
}

// ── Broadcast ─────────────────────────────────────────────────────────────────

// BroadcastPayload is the request body for broadcasting a signed transaction.
type BroadcastPayload struct {
	SignedTx  string `json:"signed_tx"`
	ChainName string `json:"chain_name,omitempty"`
	ChainID   string `json:"chain_id,omitempty"`
}

// BroadcastEVMResult is the response from an EVM broadcast.
type BroadcastEVMResult struct {
	TxHash    string `json:"tx_hash"`
	ChainName string `json:"chain_name"`
	ChainID   string `json:"chain_id"`
}

// BroadcastSolanaResult is the response from a Solana broadcast.
type BroadcastSolanaResult struct {
	Signature string `json:"signature"`
}

// ── Sweep ─────────────────────────────────────────────────────────────────────

// SweepPayload is the request body for triggering a sweep.
type SweepPayload struct {
	ChainType ChainType `json:"chain_type"`
	ChainName string    `json:"chain_name,omitempty"`
	ChainID   string    `json:"chain_id,omitempty"`
}

// ── Jobs ──────────────────────────────────────────────────────────────────────

// Job is the full state of an async operation.
type Job struct {
	ID        string         `json:"id"`
	Status    JobStatus      `json:"status"`
	Operation string         `json:"operation"`
	Result    map[string]any `json:"result,omitempty"`
	Error     string         `json:"error,omitempty"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

// ── Stablecoin ────────────────────────────────────────────────────────────────

// StablecoinTransferPayload is the request body for a stablecoin transfer.
type StablecoinTransferPayload struct {
	Token          string        `json:"token"`
	To             string        `json:"to"`
	Amount         string        `json:"amount"`
	ChainName      string        `json:"chain_name,omitempty"`
	ChainID        string        `json:"chain_id,omitempty"`
	Gasless        bool          `json:"gasless,omitempty"`
	Speed          TransferSpeed `json:"speed,omitempty"`
	IdempotencyKey string        `json:"idempotency_key,omitempty"`
}

// StablecoinTransferResult is the response from a stablecoin transfer.
type StablecoinTransferResult struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// StablecoinBalanceResult is the response from a stablecoin balance lookup.
type StablecoinBalanceResult struct {
	Address    string `json:"address"`
	Token      string `json:"token"`
	Symbol     string `json:"symbol"`
	Balance    string `json:"balance"`
	RawBalance string `json:"raw_balance"`
	ChainID    string `json:"chain_id,omitempty"`
}

// ── Withdrawals ──────────────────────────────────────────────────────────────

// CreateWithdrawalPayload is the request body for creating a withdrawal.
type CreateWithdrawalPayload struct {
	UserID         string         `json:"user_id"`
	Asset          string         `json:"asset"`
	ChainType      ChainType      `json:"chain_type"`
	ChainName      string         `json:"chain_name,omitempty"`
	ChainID        string         `json:"chain_id,omitempty"`
	Amount         string         `json:"amount"`
	To             string         `json:"to"`
	IdempotencyKey string         `json:"idempotency_key,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// ListWithdrawalsOptions controls withdrawal list filtering.
type ListWithdrawalsOptions struct {
	UserID string
	After  string
}

// Withdrawal is a treasury-safe payout request.
type Withdrawal struct {
	ID             string           `json:"id"`
	ProjectID      string           `json:"project_id,omitempty"`
	UserID         string           `json:"user_id"`
	Token          string           `json:"token"`
	ChainType      ChainType        `json:"chain_type"`
	ChainID        string           `json:"chain_id,omitempty"`
	Amount         string           `json:"amount"`
	To             string           `json:"to"`
	IdempotencyKey *string          `json:"idempotency_key,omitempty"`
	Status         WithdrawalStatus `json:"status"`
	PayoutWalletID *string          `json:"payout_wallet_id,omitempty"`
	SigningJobID   *string          `json:"signing_job_id,omitempty"`
	TxHash         *string          `json:"tx_hash,omitempty"`
	Error          *string          `json:"error,omitempty"`
	Metadata       map[string]any   `json:"metadata,omitempty"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
}

// WithdrawalList is the response from listing withdrawals.
type WithdrawalList struct {
	Withdrawals []Withdrawal `json:"withdrawals"`
}

// ── Monitor ─────────────────────────────────────────────────────────────────

// MonitorSubscriptionPayload is the request body for wallet monitoring.
//
// For EVM, set ChainType to ChainTypeEVM and pass ChainName or ChainID from
// the supported chains endpoint. For Solana and Tron, omit ChainName and
// ChainID.
type MonitorSubscriptionPayload struct {
	WalletID   string             `json:"wallet_id"`
	ChainType  ChainType          `json:"chain_type"`
	ChainName  string             `json:"chain_name,omitempty"`
	ChainID    string             `json:"chain_id,omitempty"`
	Directions []MonitorDirection `json:"directions"`
	WebhookURL string             `json:"webhook_url,omitempty"`
	MinAmount  string             `json:"min_amount,omitempty"`
	Label      string             `json:"label,omitempty"`
}

// MonitorSubscription is a wallet transaction monitoring subscription.
type MonitorSubscription struct {
	ID         string                    `json:"subscription_id"`
	WalletID   string                    `json:"wallet_id"`
	Address    string                    `json:"address"`
	ChainType  ChainType                 `json:"chain_type"`
	ChainName  string                    `json:"chain_name,omitempty"`
	ChainID    string                    `json:"chain_id,omitempty"`
	Directions []MonitorDirection        `json:"directions"`
	WebhookURL string                    `json:"webhook_url,omitempty"`
	MinAmount  string                    `json:"min_amount"`
	Label      string                    `json:"label,omitempty"`
	Status     MonitorSubscriptionStatus `json:"status"`
	CreatedAt  string                    `json:"created_at"`
}

// MonitorUnsubscribeResult is returned after removing a monitor subscription.
type MonitorUnsubscribeResult struct {
	Status string `json:"status"`
}
