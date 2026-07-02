# vaultkey-go

Official Go SDK for the [VaultKey](https://vaultkeys.dev) API.

## Installation

```bash
go get github.com/vaultkeys/vaultkey-go
```

## Requirements

- Go 1.21+
- A VaultKey API key and secret from your [dashboard](https://app.vaultkeys.dev)

---

## Quick Start

```go
package main

import (
    "context"
    "log"

    vaultkey "github.com/vaultkey/vaultkey-go"
)

func main() {
    client, err := vaultkey.NewClient("vk_live_...", "your_api_secret")
    if err != nil {
        log.Fatal(err)
    }

    wallet, apiErr, err := client.Wallets.Create(context.Background(), vaultkey.CreateWalletPayload{
        UserID:    "user_123",
        ChainType: vaultkey.ChainTypeEVM,
    })
    if err != nil {
        log.Fatal(err)
    }
    if apiErr != nil {
        log.Fatalf("api error: %s", apiErr.Message)
    }

    log.Printf("wallet created: %s %s", wallet.ID, wallet.Address)
}
```

API key and secret can also be supplied via environment variables:

```bash
export VAULTKEY_API_KEY=vk_live_...
export VAULTKEY_API_SECRET=...
```

---

## Configuration

```go
client, err := vaultkey.NewClient(
    "vk_live_...",
    "your_api_secret",
    vaultkey.WithBaseURL("https://your-self-hosted-instance.com/api/v1/sdk"),
    vaultkey.WithHTTPClient(myHTTPClient),
)
```

**Options:**
- `WithBaseURL(url)` — override the base URL for self-hosted deployments. When set, takes precedence over automatic key-based routing.
- `WithHTTPClient(*http.Client)` — provide a custom HTTP client (timeouts, proxies, etc.).

**Key prefix routing (automatic):**
- `testnet_` keys → `https://testnet.vaultkeys.dev`
- `vk_live_` keys → `https://app.vaultkeys.dev`

**Default HTTP client** uses a 30s timeout.

---

## Error Handling

Every method returns `(result, *ErrorResponse, error)`:
- `error` — transport/network error (no response received)
- `*ErrorResponse` — API returned a non-2xx response
- Both are nil on success

```go
wallet, apiErr, err := client.Wallets.Get(ctx, "wallet_id")
if err != nil {
    // network or serialization error
    log.Fatal(err)
}
if apiErr != nil {
    // API-level error (404, 403, etc.)
    log.Fatalf("%s: %s", apiErr.Code, apiErr.Message)
}
log.Println(wallet.Address)
```

---

## Wallets

### Create

```go
wallet, apiErr, err := client.Wallets.Create(ctx, vaultkey.CreateWalletPayload{
    UserID:    "user_123",
    ChainType: vaultkey.ChainTypeEVM,   // ChainTypeEVM | ChainTypeSolana | ChainTypeTron
    Label:     "Primary",               // optional
})
```

### Get

```go
wallet, apiErr, err := client.Wallets.Get(ctx, "wallet_id")
```

### List by user

```go
result, apiErr, err := client.Wallets.ListByUser(ctx, "user_123", "", 20)

// Next page
if result.HasMore {
    page2, _, _ := client.Wallets.ListByUser(ctx, "user_123", result.NextCursor, 20)
}
```

---

## Signing

Signing operations are **asynchronous** — they return a job ID which you poll via `Jobs.Get`.

```go
signing := client.Wallets.Signing("wallet_id")

// EVM message
job, apiErr, err := signing.EVMMessage(ctx, vaultkey.SignMessagePayload{
    Payload:        map[string]any{"message": "Hello from VaultKey"},
    IdempotencyKey: "sign-001", // optional — safe to retry
})

// Solana message
job, apiErr, err := signing.SolanaMessage(ctx, vaultkey.SignMessagePayload{
    Payload: map[string]any{"data": "SGVsbG8="},
})

// Poll until done
result, apiErr, err := client.Jobs.Get(ctx, job.JobID)
```

---

## Balances

```go
// EVM — chainName (preferred) or chainID
bal, apiErr, err := client.Wallets.EVMBalance(ctx, "wallet_id", "base", "")
fmt.Println(bal.Balance, bal.Symbol) // "0.05" "ETH"

// Solana
bal, apiErr, err := client.Wallets.SolanaBalance(ctx, "wallet_id")
fmt.Println(bal.Balance) // "1.5"
```

---

## Broadcast

```go
// EVM
result, apiErr, err := client.Wallets.BroadcastEVM(ctx, "wallet_id", "0x...", "base", "")
fmt.Println(result.TxHash)

// Solana
result, apiErr, err := client.Wallets.BroadcastSolana(ctx, "wallet_id", "base58tx...")
fmt.Println(result.Signature)
```

---

## Sweep

```go
// EVM sweep
job, apiErr, err := client.Wallets.Sweep(ctx, "wallet_id", vaultkey.SweepPayload{
    ChainType: vaultkey.ChainTypeEVM,
    ChainName: "base",
})

// Solana sweep
job, apiErr, err := client.Wallets.Sweep(ctx, "wallet_id", vaultkey.SweepPayload{
    ChainType: vaultkey.ChainTypeSolana,
})

result, _, _ := client.Jobs.Get(ctx, job.JobID)
```

---

## Withdrawals

Use withdrawals for outbound stablecoin payouts. VaultKey selects the payout
wallet and applies treasury policy behind the API.

```go
withdrawal, apiErr, err := client.Withdrawals.Create(ctx, vaultkey.CreateWithdrawalPayload{
    UserID:         "user_123",
    Asset:          "usdc",
    ChainType:      vaultkey.ChainTypeEVM,
    ChainName:      "base",
    Amount:         "100.00",
    To:             "0x1111111111111111111111111111111111111111",
    IdempotencyKey: "withdrawal_user_123_001",
    Metadata: map[string]any{
        "external_reference": "wd_001",
    },
})

withdrawal, apiErr, err = client.Withdrawals.Get(ctx, withdrawal.ID)

list, apiErr, err := client.Withdrawals.List(ctx, vaultkey.ListWithdrawalsOptions{
    UserID: "user_123",
})

withdrawal, apiErr, err = client.Withdrawals.Cancel(ctx, withdrawal.ID)
```

Common statuses include `pending_approval`, `funding_required`, `funding`,
`signing`, `completed`, `failed`, and `cancelled`.

---

## Monitor

Subscribe wallets to incoming and/or outgoing transaction webhooks.

```go
sub, apiErr, err := client.Monitor.Subscribe(ctx, vaultkey.MonitorSubscriptionPayload{
    WalletID:  "wallet_id",
    ChainType: vaultkey.ChainTypeEVM,
    ChainName: "sepolia",
    // ChainID: "11155111", // optional fallback if you do not use ChainName
    Directions: []vaultkey.MonitorDirection{
        vaultkey.MonitorDirectionIncoming,
        vaultkey.MonitorDirectionOutgoing,
    },
    WebhookURL: "https://example.com/webhook",
    MinAmount:  "0",
    Label:      "optional",
})

subs, apiErr, err := client.Monitor.List(ctx)

sub, apiErr, err = client.Monitor.Get(ctx, sub.ID)

result, apiErr, err := client.Monitor.Unsubscribe(ctx, sub.ID)
fmt.Println(result.Status)
```

---

## Stablecoin

Direct stablecoin transfers are a lower-level wallet operation. For normal
outbound payouts, prefer `client.Withdrawals`.

```go
// Transfer USDC on Base (gasless)
result, apiErr, err := client.Stablecoin.Transfer(ctx, "wallet_id",
    vaultkey.ChainTypeEVM,
    vaultkey.StablecoinTransferPayload{
        Token:          "usdc",
        To:             "0xRecipient",
        Amount:         "50.00",
        ChainName:      "base",
        Gasless:        true,
        Speed:          vaultkey.SpeedFast,
        IdempotencyKey: "tx-001", // optional — prevents double sends
    },
)

// Solana transfer — omit chain fields
result, apiErr, err := client.Stablecoin.Transfer(ctx, "wallet_id",
    vaultkey.ChainTypeSolana,
    vaultkey.StablecoinTransferPayload{
        Token:  "usdc",
        To:     "RecipientBase58...",
        Amount: "50.00",
    },
)

// Poll
job, _, _ := client.Jobs.Get(ctx, result.JobID)

// Balance
bal, apiErr, err := client.Stablecoin.Balance(ctx, "wallet_id",
    vaultkey.ChainTypeEVM, "usdc", "polygon", "",
)
fmt.Println(bal.Balance) // "50.00"
```

---

## Jobs

```go
result, apiErr, err := client.Jobs.Get(ctx, "job_id")
// result.Status: JobStatusPending | JobStatusProcessing | JobStatusCompleted | JobStatusFailed
```

### Polling helper

```go
func pollJob(ctx context.Context, client *vaultkey.Client, jobID string) (vaultkey.Job, error) {
    for {
        result, apiErr, err := client.Jobs.Get(ctx, jobID)
        if err != nil {
            return vaultkey.Job{}, err
        }
        if apiErr != nil {
            return vaultkey.Job{}, fmt.Errorf("%s: %s", apiErr.Code, apiErr.Message)
        }
        if result.Status == vaultkey.JobStatusCompleted {
            return result, nil
        }
        if result.Status == vaultkey.JobStatusFailed {
            return result, fmt.Errorf("job failed: %s", result.Error)
        }
        time.Sleep(time.Second)
    }
}
```

---

## Chains

```go
chains, apiErr, err := client.Chains.List(ctx)
for _, c := range chains {
    fmt.Printf("%s  chain_id=%s  %s\n", c.Name, c.ChainID, c.NativeSymbol)
}
```

---

## Available Services

| Service | Methods |
|---|---|
| `client.Wallets` | `Create`, `Get`, `ListByUser`, `EVMBalance`, `SolanaBalance`, `BroadcastEVM`, `BroadcastSolana`, `Sweep`, `Signing(id)` |
| `client.Wallets.Signing(id)` | `EVMMessage`, `SolanaMessage` |
| `client.Withdrawals` | `Create`, `Get`, `List`, `Cancel` |
| `client.Monitor` | `Subscribe`, `List`, `Get`, `Unsubscribe` |
| `client.Stablecoin` | `Transfer`, `Balance` |
| `client.Jobs` | `Get` |
| `client.Chains` | `List` |
