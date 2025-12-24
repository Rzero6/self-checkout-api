package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go/coreapi"
)

func VerifyMidtransSignature(notif coreapi.TransactionStatusResponse, serverKey string) bool {
	// signature string = order_id + status_code + gross_amount + serverKey
	str := fmt.Sprintf("%s%s%s%s",
		notif.OrderID,
		notif.StatusCode,
		notif.GrossAmount,
		serverKey,
	)

	hash := sha512.Sum512([]byte(str))
	computedSignature := hex.EncodeToString(hash[:])

	return computedSignature == notif.SignatureKey
}

func TranslateMidtransPaymentStatus(midtransStatus coreapi.TransactionStatusResponse) string {
	var localStatus string
	switch midtransStatus.TransactionStatus {
	case "capture":
		if midtransStatus.FraudStatus == "accept" {
			localStatus = string(TransactionStatusSuccess)
		} else {
			localStatus = string(TransactionStatusFailed)
		}
	case "settlement":
		localStatus = string(TransactionStatusSuccess)
	case "pending":
		localStatus = string(TransactionStatusPending)
	case "expire":
		localStatus = string(TransactionStatusExpired)
	case "cancel":
		localStatus = string(TransactionStatusCancelled)
	case "deny", "failure":
		localStatus = string(TransactionStatusFailed)
	default:
		return strings.ToUpper(midtransStatus.TransactionStatus)
	}
	return localStatus
}

func ParseJakartaToUTC(timeStr string) (*time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("failed to load Jakarta location: %w", err)
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	utcTime := t.UTC()
	return &utcTime, nil
}
