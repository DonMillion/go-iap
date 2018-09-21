package appstore

import "encoding/json"

type numericString string

func (n *numericString) UnmarshalJSON(b []byte) error {
	var number json.Number
	if err := json.Unmarshal(b, &number); err != nil {
		return err
	}
	*n = numericString(number.String())
	return nil
}

type Environment string

const (
	Sandbox    Environment = "Sandbox"
	Production Environment = "Production"
)

type (
	// https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
	// The IAPRequest type has the request parameter
	IAPRequest struct {
		ReceiptData string `json:"receipt-data"`
		// Only used for receipts that contain auto-renewable subscriptions.
		Password string `json:"password,omitempty"`
		// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
		// If value is true, response includes only the latest renewal transaction for any subscriptions.
		ExcludeOldTransactions bool `json:"exclude-old-transactions"`
	}

	// The ReceiptCreationDate type indicates the date when the app receipt was created.
	ReceiptCreationDate struct {
		CreationDate    string `json:"receipt_creation_date"`
		CreationDateMS  string `json:"receipt_creation_date_ms"`
		CreationDatePST string `json:"receipt_creation_date_pst"`
	}

	// The RequestDate type indicates the date and time that the request was sent
	RequestDate struct {
		RequestDate    string `json:"request_date"`
		RequestDateMS  string `json:"request_date_ms"`
		RequestDatePST string `json:"request_date_pst"`
	}

	// The PurchaseDate type indicates the date and time that the item was purchased
	PurchaseDate struct {
		PurchaseDate    string `json:"purchase_date"`
		PurchaseDateMS  string `json:"purchase_date_ms"`
		PurchaseDatePST string `json:"purchase_date_pst"`
	}

	// The OriginalPurchaseDate type indicates the beginning of the subscription period
	OriginalPurchaseDate struct {
		OriginalPurchaseDate    string `json:"original_purchase_date"`
		OriginalPurchaseDateMS  string `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePST string `json:"original_purchase_date_pst"`
	}

	// The ExpiresDate type indicates the expiration date for the subscription
	ExpiresDate struct {
		ExpiresDate             string `json:"expires_date,omitempty"`
		ExpiresDateMS           string `json:"expires_date_ms,omitempty"`
		ExpiresDatePST          string `json:"expires_date_pst,omitempty"`
		ExpiresDateFormatted    string `json:"expires_date_formatted,omitempty"`
		ExpiresDateFormattedPST string `json:"expires_date_formatted_pst,omitempty"`
	}

	// The CancellationDate type indicates the time and date of the cancellation by Apple customer support
	CancellationDate struct {
		CancellationDate    string `json:"cancellation_date,omitempty"`
		CancellationDateMS  string `json:"cancellation_date_ms,omitempty"`
		CancellationDatePST string `json:"cancellation_date_pst,omitempty"`
	}

	// The InApp type has the receipt attributes
	InApp struct {
		Quantity              string `json:"quantity"`
		ProductID             string `json:"product_id"`
		TransactionID         string `json:"transaction_id"`
		OriginalTransactionID string `json:"original_transaction_id"`
		WebOrderLineItemID    string `json:"web_order_line_item_id,omitempty"`

		IsTrialPeriod string `json:"is_trial_period"`
		ExpiresDate

		PurchaseDate
		OriginalPurchaseDate

		CancellationDate
		CancellationReason string `json:"cancellation_reason,omitempty"`
	}

	// The Receipt type has whole data of receipt
	Receipt struct {
		ReceiptType                string        `json:"receipt_type"`
		AdamID                     int64         `json:"adam_id"`
		AppItemID                  numericString `json:"app_item_id"`
		BundleID                   string        `json:"bundle_id"`
		ApplicationVersion         string        `json:"application_version"`
		DownloadID                 int64         `json:"download_id"`
		VersionExternalIdentifier  numericString `json:"version_external_identifier"`
		OriginalApplicationVersion string        `json:"original_application_version"`
		InApp                      []InApp       `json:"in_app"`
		ReceiptCreationDate
		RequestDate
		OriginalPurchaseDate
	}

	// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
	PendingRenewalInfo struct {
		SubscriptionExpirationIntent   string `json:"expiration_intent"`
		SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
		SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
		SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
		SubscriptionPriceConsentStatus string `json:"price_consent_status"`
		ProductID                      string `json:"product_id"`
	}

	// The IAPResponse type has the response properties
	// We defined each field by the current IAP response, but some fields are not mentioned
	// in the following Apple's document;
	// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html
	// If you get other types or fields from the IAP response, you should use the struct you defined.
	IAPResponse struct {
		Status             int                  `json:"status"`
		Environment        Environment          `json:"environment"`
		Receipt            Receipt              `json:"receipt"`
		LatestReceiptInfo  []InApp              `json:"latest_receipt_info,omitempty"`
		LatestReceipt      string               `json:"latest_receipt,omitempty"`
		PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info,omitempty"`
		IsRetryable        bool                 `json:"is-retryable,omitempty"`
	}

	// The HttpStatusResponse struct contains the status code returned by the store
	// Used as a workaround to detect when to hit the production appstore or sandbox appstore regardless of receipt type
	StatusResponse struct {
		Status int `json:"status"`
	}

	// IAPResponseForIOS6 is iOS 6 style receipt schema.
	IAPResponseForIOS6 struct {
		AutoRenewProductID     string         `json:"auto_renew_product_id"`
		AutoRenewStatus        int            `json:"auto_renew_status"`
		CancellationReason     string         `json:"cancellation_reason,omitempty"`
		ExpirationIntent       string         `json:"expiration_intent,omitempty"`
		IsInBillingRetryPeriod string         `json:"is_in_billing_retry_period,omitempty"`
		LatestReceiptInfo      ReceiptForIOS6 `json:"latest_expired_receipt_info"`
		Receipt                ReceiptForIOS6 `json:"receipt"`
		Status                 int            `json:"status"`
	}

	ReceiptForIOS6 struct {
		AppItemID numericString `json:"app_item_id"`
		BID       string        `json:"bid"`
		BVRS      string        `json:"bvrs"`
		CancellationDate
		ExpiresDate
		IsTrialPeriod        string `json:"is_trial_period"`
		IsInIntroOfferPeriod string `json:"is_in_intro_offer_period"`
		ItemID               string `json:"item_id"`
		ProductID            string `json:"product_id"`
		PurchaseDate
		OriginalTransactionID string `json:"original_transaction_id"`
		OriginalPurchaseDate
		Quantity                  string        `json:"quantity"`
		TransactionID             string        `json:"transaction_id"`
		UniqueIdentifier          string        `json:"unique_identifier"`
		UniqueVendorIdentifier    string        `json:"unique_vendor_identifier"`
		VersionExternalIdentifier numericString `json:"version_external_identifier,omitempty"`
		WebOrderLineItemID        string        `json:"web_order_line_item_id"`
	}

	// 购买返回的receipt (receipt for purchase result)
	PurchaseReceipt struct {
		Quantity                  string `json:"quantity,omitempty"`                    // "quantity": "1",
		UniqueVendorIdentifier    string `json:"unique_vendor_identifier,omitempty"`    // "unique_vendor_identifier": "C8C17B66-394D-46C6-996C-3AF7CF46876C",
		Bvrs                      string `json:"bvrs,omitempty"`                        // "bvrs": "9",
		AppItemId                 string `json:"app_item_id,omitempty"`                 // "app_item_id": "557130558",
		ExpiresDate               string `json:"expires_date,omitempty"`                // "expires_date": "1538917472000",
		ExpiresDateFormatted      string `json:"expires_date_formatted,omitempty"`      // "expires_date_formatted": "2018-10-07 13:04:32 Etc/GMT",
		ExpiresDateFormattedPST   string `json:"expires_date_formatted_pst,omitempty"`  // "expires_date_formatted_pst": "2018-10-07 06:04:32 America/Los_Angeles",
		IsInIntroOfferPeriod      string `json:"is_in_intro_offer_period,omitempty"`    // "is_in_intro_offer_period": "false",
		IsTrialPeriod             string `json:"is_trial_period,omitempty"`             // "is_trial_period": "false",
		ItemId                    string `json:"item_id,omitempty"`                     // "item_id": "1141778299",
		UniqueIdentifier          string `json:"unique_identifier,omitempty"`           // "unique_identifier": "e0d914df721ba7e321222382732ed4f38962b7c5",
		OriginalTransactionId     string `json:"original_transaction_id,omitempty"`     // "original_transaction_id": "220000492510109",
		TransactionId             string `json:"transaction_id,omitempty"`              // "transaction_id": "220000492510109",
		WebOrderLineItemId        string `json:"web_order_line_item_id,omitempty"`      // "web_order_line_item_id": "220000131547633",
		Bid                       string `json:"bid,omitempty"`                         // "bid": "com.helloTalk.helloTalk",
		ProductId                 string `json:"product_id,omitempty"`                  // "product_id": "com.hellotalk.monthauto",
		PurchaseDate              string `json:"purchase_date,omitempty"`               // "purchase_date": "2018-09-07 13:04:32 Etc/GMT",
		PurchaseDateMS            string `json:"purchase_date_ms,omitempty"`            // "purchase_date_ms": "1536325472000",
		PurchaseDatePST           string `json:"purchase_date_pst,omitempty"`           // "purchase_date_pst": "2018-09-07 06:04:32 America/Los_Angeles",
		OriginalPurchaseDate      string `json:"original_purchase_date,omitempty"`      // "original_purchase_date": "2018-09-07 13:04:34 Etc/GMT",
		OriginalPurchaseDateMS    string `json:"original_purchase_date_ms,omitempty"`   // "original_purchase_date_ms": "1536325474000",
		OriginalPurchaseDatePST   string `json:"original_purchase_date_pst,omitempty"`  // "original_purchase_date_pst": "2018-09-07 06:04:34 America/Los_Angeles",
		VersionExternalIdentifier string `json:"version_external_identifier,omitempty"` // "version_external_identifier": "828156491",
		BundleID                  string `json:"bundle_id"`
		ApplicationVersion        string `json:"application_version"`
		CreationDate              string `json:"receipt_creation_date"`
		CreationDateMS            string `json:"receipt_creation_date_ms"`
		CreationDatePST           string `json:"receipt_creation_date_pst"`
		IsInBillingRetryPeriod    string `is_in_billing_retry_period`
	}

	PurchaseIAPResponse struct {
		AutoRenewStatus         uint32          `json:"auto_renew_status"` // 续费商品才有
		Status                  uint32          `json:"status"`
		AutoRenewProductId      string          `json:"auto_renew_product_id"` // 续费商品才有
		Receipt                 PurchaseReceipt `json:"receipt"`
		LatestReceiptInfo       PurchaseReceipt `json:"latest_receipt_info"`         // 续费商品才有
		LatestExpireReceiptInfo PurchaseReceipt `json:"latest_expired_receipt_info"` // 续费商品才有
		LatestReceipt           string          `json:"latest_receipt"`              // 续费商品才有
	}
)
