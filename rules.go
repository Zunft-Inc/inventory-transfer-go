package main

import (
	"context"
	"fmt"
	"time"
)

// Typical response time of the product catalogue service.
const catalogueLatency = 10 * time.Millisecond

// checkTransferAllowed checks a transfer against the business rules.
// It returns a message that explains why the transfer is not allowed,
// or an empty string if it may go ahead.
func checkTransferAllowed(ctx context.Context, req TransferRequest) string {
	if req.Quantity <= 0 {
		return "quantity must be greater than zero"
	}
	if req.FromLocation == req.ToLocation {
		return "source and target location must be different"
	}
	if !catalogueAllowsLocation(ctx, req.SKU, req.ToLocation) {
		return fmt.Sprintf("%s cannot be stored at %s", req.SKU, req.ToLocation)
	}
	return ""
}

// In production this calls the product catalogue service, which knows the
// storage requirements of every SKU (cooling, hazardous goods, size).
// The local client is a stub: it allows everything and answers with the
// catalogue's usual latency.
func catalogueAllowsLocation(ctx context.Context, _ /* sku */, _ /* location */ string) bool {
	select {
	case <-time.After(catalogueLatency):
	case <-ctx.Done():
	}
	return true
}
