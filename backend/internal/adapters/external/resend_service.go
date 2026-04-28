package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend/internal/adapters/logging"
	"backend/internal/core/ports"
)

type ResendService struct {
	apiKey string
	client *http.Client
	log    ports.Logger
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
		log: logging.With("email"),
	}
}

// SendWelcomeEmail sends a welcome email to a new beta user with activation link
func (r *ResendService) SendWelcomeEmail(email, activationLink string) error {
	r.log.Debug("Sending welcome email",
		ports.F("email", email),
		ports.F("has_activation_link", activationLink != ""),
	)

	if r.apiKey == "" {
		r.log.Error("RESEND_API_KEY not configured")
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

// SendDeletionConfirmationEmail sends a confirmation email when account deletion is requested
func (r *ResendService) SendDeletionConfirmationEmail(email, cancelURL string, deletesAt time.Time) error {
	r.log.Debug("Sending deletion confirmation email",
		ports.F("email", email),
		ports.F("deletes_at", deletesAt),
	)

	if r.apiKey == "" {
		r.log.Error("RESEND_API_KEY not configured")
		return fmt.Errorf("RESEND_API_KEY not configured")
	}

	deletionEmail := ResendEmail{
		From:    "CannaNote <noreply@mail.cannanote.org>",
		To:      []string{email},
		Subject: "Confirm Account Deletion - CannaNote",
		HTML: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Account Deletion Confirmation</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #dc2626 0%%, #991b1b 100%%); color: white; padding: 30px; border-radius: 8px 8px 0 0; text-align: center; }
        .content { background: white; padding: 30px; border-radius: 0 0 8px 8px; border: 1px solid #e2e8f0; }
        .cta { background: #48bb78; color: white; padding: 16px 32px; text-decoration: none; border-radius: 6px; display: inline-block; margin: 20px 0; font-weight: bold; font-size: 16px; }
        .cta:hover { background: #38a169; }
        .footer { text-align: center; margin-top: 30px; color: #718096; font-size: 14px; }
        .warning { background: #fef3cd; padding: 15px; border-radius: 6px; margin: 20px 0; border-left: 4px solid #ffc107; }
        .countdown { font-size: 24px; font-weight: bold; color: #dc2626; text-align: center; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Account Deletion Requested</h1>
        </div>
        <div class="content">
            <p>Hi there,</p>
            <p>You've requested to delete your CannaNote account. We're sorry to see you go.</p>

            <div class="warning">
                <strong>Important:</strong> Your account and all data will be permanently deleted on:
                <div class="countdown">%s</div>
            </div>

            <p>During this 24-hour grace period, you can:</p>
            <ul>
                <li>Continue using CannaNote normally</li>
                <li>Export your data from Settings</li>
                <li>Cancel the deletion if you change your mind</li>
            </ul>

            <h3>Didn't request this?</h3>
            <p>If you didn't request to delete your account, click the button below to cancel immediately:</p>

            <center>
                <a href="%s" class="cta">Cancel Deletion</a>
            </center>

            <p style="font-size: 13px; color: #666;">This link expires when the deletion is processed. You can also cancel from within the app.</p>

            <h3>What happens after deletion?</h3>
            <ul>
                <li>All your session history will be erased</li>
                <li>Your account credentials will be removed</li>
                <li>Local data on your devices will be cleared</li>
                <li>This action cannot be undone</li>
            </ul>

            <p>If you're leaving because something wasn't working right, we'd love to hear from you. Reply to this email with any feedback.</p>

            <p>Thank you for being part of the CannaNote community.</p>
            <p>The CannaNote Team</p>
        </div>
        <div class="footer">
            <p>CannaNote - Privacy-focused cannabis tracking<br>
            <a href="https://cannanote.org/privacy">Privacy Policy</a> | <a href="https://cannanote.org/terms">Terms of Service</a></p>
        </div>
    </div>
</body>
</html>
		`, deletesAt.Format("January 2, 2006 at 3:04 PM MST"), cancelURL),
		Text: fmt.Sprintf(`Account Deletion Requested

Hi there,

You've requested to delete your CannaNote account. We're sorry to see you go.

IMPORTANT: Your account and all data will be permanently deleted on:
%s

During this 24-hour grace period, you can:
- Continue using CannaNote normally
- Export your data from Settings
- Cancel the deletion if you change your mind

DIDN'T REQUEST THIS?
If you didn't request to delete your account, cancel immediately:
%s

This link expires when the deletion is processed. You can also cancel from within the app.

WHAT HAPPENS AFTER DELETION?
- All your session history will be erased
- Your account credentials will be removed
- Local data on your devices will be cleared
- This action cannot be undone

If you're leaving because something wasn't working right, we'd love to hear from you. Reply to this email with any feedback.

Thank you for being part of the CannaNote community.
The CannaNote Team

---
CannaNote - Privacy-focused cannabis tracking
Privacy Policy: https://cannanote.org/privacy
Terms of Service: https://cannanote.org/terms`, deletesAt.Format("January 2, 2006 at 3:04 PM MST"), cancelURL),
	}

	return r.sendEmail(deletionEmail)
}

// sendEmail sends an email via Resend API
func (r *ResendService) sendEmail(email ResendEmail) error {
	r.log.Debug("Preparing email",
		ports.F("to", email.To),
		ports.F("subject", email.Subject),
	)

	jsonData, err := json.Marshal(email)
	if err != nil {
		r.log.Error("Failed to marshal email data", ports.F("error", err))
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		r.log.Error("Failed to create HTTP request", ports.F("error", err))
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		r.log.Error("HTTP request failed", ports.F("error", err))
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp ResendResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			r.log.Error("Resend API error",
				ports.F("status", resp.StatusCode),
				ports.F("message", errorResp.Message),
			)
			return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		r.log.Error("Resend API error", ports.F("status", resp.StatusCode))
		return fmt.Errorf("resend API error: %d", resp.StatusCode)
	}

	var successResp ResendResponse
	if err := json.NewDecoder(resp.Body).Decode(&successResp); err == nil {
		r.log.Info("Email sent successfully",
			ports.F("email_id", successResp.ID),
			ports.F("to", email.To),
		)
	} else {
		r.log.Info("Email sent successfully", ports.F("to", email.To))
	}

	return nil
}