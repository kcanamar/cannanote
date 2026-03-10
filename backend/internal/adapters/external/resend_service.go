package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ResendService struct {
	apiKey string
	client *http.Client
}

type ResendEmail struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	HTML     string   `json:"html"`
	Text     string   `json:"text,omitempty"`
}

type ResendResponse struct {
	ID      string `json:"id"`
	Message string `json:"message,omitempty"`
}

// NewResendService creates a new Resend email service
func NewResendService() *ResendService {
	return &ResendService{
		apiKey: os.Getenv("RESEND_API_KEY"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendWelcomeEmail sends a welcome email to a new beta user with activation link
func (r *ResendService) SendWelcomeEmail(email, activationLink string) error {
	fmt.Printf("[DEBUG] SendWelcomeEmail: Starting for email=%s\n", email)
	fmt.Printf("[DEBUG] SendWelcomeEmail: activationLink=%s\n", activationLink)
	fmt.Printf("[DEBUG] SendWelcomeEmail: RESEND_API_KEY configured: %v (length: %d)\n", r.apiKey != "", len(r.apiKey))

	if r.apiKey == "" {
		fmt.Printf("[DEBUG] SendWelcomeEmail: FAILED - RESEND_API_KEY not configured\n")
		return fmt.Errorf("RESEND_API_KEY not configured")
	}

	// Fallback if no activation link provided
	if activationLink == "" {
		activationLink = "https://cannanote.org"
	}

	welcomeEmail := ResendEmail{
		From:    "CannaNote <welcome@mail.cannanote.org>",
		To:      []string{email},
		Subject: "Activate Your CannaNote Beta Account",
		HTML: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Activate Your CannaNote Account</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #1a1f36 0%%, #2d3748 100%%); color: white; padding: 30px; border-radius: 8px 8px 0 0; text-align: center; }
        .content { background: white; padding: 30px; border-radius: 0 0 8px 8px; border: 1px solid #e2e8f0; }
        .cta { background: #48bb78; color: white; padding: 16px 32px; text-decoration: none; border-radius: 6px; display: inline-block; margin: 20px 0; font-weight: bold; font-size: 16px; }
        .cta:hover { background: #38a169; }
        .footer { text-align: center; margin-top: 30px; color: #718096; font-size: 14px; }
        .note { background: #f7fafc; padding: 15px; border-radius: 6px; margin: 20px 0; border-left: 4px solid #48bb78; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to CannaNote</h1>
            <p>Your Personal Cannabis Wellness Companion</p>
        </div>
        <div class="content">
            <p>Hi there,</p>
            <p>Welcome to the CannaNote beta! You've just joined a community focused on mindful cannabis consumption and personal wellness insights.</p>

            <div class="note">
                <strong>One click to get started:</strong> Click the button below to activate your account and set your password.
            </div>

            <center>
                <a href="%s" class="cta">Activate My Account</a>
            </center>

            <h3>What You'll Get</h3>
            <ul>
                <li><strong>30-second logging</strong> - Track sessions quickly on any device</li>
                <li><strong>Pattern insights</strong> - Discover what works best for you</li>
                <li><strong>Privacy first</strong> - Your data stays local unless you choose sync</li>
                <li><strong>Beta perks</strong> - Lifetime premium access for early supporters</li>
            </ul>

            <h3>Beta Grandfathering</h3>
            <p>As a beta member, you'll receive <strong>lifetime access</strong> to premium sync features once we launch. Your early support means everything to us!</p>

            <h3>Need Help?</h3>
            <p>Visit our <a href="https://cannanote.org/docs">documentation</a> or reply to this email with any questions.</p>

            <p>Thank you for joining our cannabis wellness journey!</p>
            <p>The CannaNote Team</p>
        </div>
        <div class="footer">
            <p>CannaNote - Privacy-focused cannabis tracking<br>
            <a href="https://cannanote.org/privacy">Privacy Policy</a> | <a href="https://cannanote.org/terms">Terms of Service</a></p>
            <p style="font-size: 12px; color: #a0aec0;">If you didn't sign up for CannaNote, you can safely ignore this email.</p>
        </div>
    </div>
</body>
</html>
		`, activationLink),
		Text: fmt.Sprintf(`Welcome to CannaNote Beta!

Hi there,

Welcome to the CannaNote beta! You've just joined a community focused on mindful cannabis consumption and personal wellness insights.

ACTIVATE YOUR ACCOUNT
Click the link below to verify your email and set your password:
%s

WHAT YOU'LL GET
- 30-second logging - Track sessions quickly on any device
- Pattern insights - Discover what works best for you
- Privacy first - Your data stays local unless you choose sync
- Beta perks - Lifetime premium access for early supporters

BETA GRANDFATHERING
As a beta member, you'll receive lifetime access to premium sync features once we launch. Your early support means everything to us!

NEED HELP?
Visit our documentation at https://cannanote.org/docs or reply to this email.

Thank you for joining our cannabis wellness journey!
The CannaNote Team

---
CannaNote - Privacy-focused cannabis tracking
Privacy Policy: https://cannanote.org/privacy
Terms of Service: https://cannanote.org/terms

If you didn't sign up for CannaNote, you can safely ignore this email.`, activationLink),
	}

	return r.sendEmail(welcomeEmail)
}

// sendEmail sends an email via Resend API
func (r *ResendService) sendEmail(email ResendEmail) error {
	fmt.Printf("[DEBUG] sendEmail: Preparing to send email to %v\n", email.To)
	fmt.Printf("[DEBUG] sendEmail: From=%s, Subject=%s\n", email.From, email.Subject)

	jsonData, err := json.Marshal(email)
	if err != nil {
		fmt.Printf("[DEBUG] sendEmail: FAILED to marshal email data: %v\n", err)
		return fmt.Errorf("failed to marshal email data: %w", err)
	}
	fmt.Printf("[DEBUG] sendEmail: Email JSON marshaled successfully (size: %d bytes)\n", len(jsonData))

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[DEBUG] sendEmail: FAILED to create HTTP request: %v\n", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")
	fmt.Printf("[DEBUG] sendEmail: Making POST request to Resend API...\n")

	resp, err := r.client.Do(req)
	if err != nil {
		fmt.Printf("[DEBUG] sendEmail: FAILED HTTP request: %v\n", err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("[DEBUG] sendEmail: Resend API response status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		var errorResp ResendResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			fmt.Printf("[DEBUG] sendEmail: FAILED - Resend API error (%d): %s\n", resp.StatusCode, errorResp.Message)
			return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		fmt.Printf("[DEBUG] sendEmail: FAILED - Resend API error: %d\n", resp.StatusCode)
		return fmt.Errorf("resend API error: %d", resp.StatusCode)
	}

	var successResp ResendResponse
	if err := json.NewDecoder(resp.Body).Decode(&successResp); err == nil {
		fmt.Printf("[DEBUG] sendEmail: SUCCESS - Email sent with ID: %s\n", successResp.ID)
	} else {
		fmt.Printf("[DEBUG] sendEmail: SUCCESS - Email sent (no ID in response)\n")
	}

	return nil
}