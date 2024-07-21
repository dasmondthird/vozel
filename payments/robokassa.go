package payments

import (
        "crypto/md5"
        "fmt"
        "net/url"
        "os"
        "strconv"
)

type RoboKassa struct {
        MerchantLogin string
        Password1     string
        Password2     string
        TestMode      bool
}

func NewRoboKassa() *RoboKassa {
        return &RoboKassa{
                MerchantLogin: os.Getenv("ROBOKASSA_MERCHANT_LOGIN"),
                Password1:     os.Getenv("ROBOKASSA_PASSWORD1"),
                Password2:     os.Getenv("ROBOKASSA_PASSWORD2"),
                TestMode:      os.Getenv("ROBOKASSA_TEST_MODE") == "true",
        }
}

func (r *RoboKassa) GeneratePaymentURL(invoiceID int, amount int, description string) string {
        baseURL := "https://auth.robokassa.ru/Merchant/Index.aspx"
        if r.TestMode {
                baseURL = "https://auth.robokassa.ru/Merchant/PaymentForm/FormTestFs.aspx"
        }

        params := url.Values{}
        params.Add("MerchantLogin", r.MerchantLogin)
        params.Add("OutSum", strconv.Itoa(amount))
        params.Add("InvId", strconv.Itoa(invoiceID))
        params.Add("Description", description)
        params.Add("SignatureValue", r.generateSignature(amount, invoiceID))
        params.Add("IsTest", strconv.Itoa(map[bool]int{true: 1, false: 0}[r.TestMode]))

        return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func (r *RoboKassa) generateSignature(amount int, invoiceID int) string {
        signature := fmt.Sprintf("%s:%d:%s:%d", r.MerchantLogin, amount, r.Password1, invoiceID)
        return fmt.Sprintf("%x", md5.Sum([]byte(signature)))
}