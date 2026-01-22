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

// SendWelcomeEmail sends a welcome email to a new beta user
func (r *ResendService) SendWelcomeEmail(email, name string) error {
	if r.apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY not configured")
	}

	welcomeEmail := ResendEmail{
		From:    "CannaNote <welcome@mail.cannanote.org>",
		To:      []string{email},
		Subject: "Welcome to CannaNote Beta",
		HTML: fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Welcome to CannaNote Beta</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #1a1f36 0%%, #2d3748 100%%); color: white; padding: 30px; border-radius: 8px 8px 0 0; text-align: center; }
        .content { background: white; padding: 30px; border-radius: 0 0 8px 8px; border: 1px solid #e2e8f0; }
        .cta { background: #48bb78; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; margin: 20px 0; }
        .footer { text-align: center; margin-top: 30px; color: #718096; font-size: 14px; }
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
            
            <h3>What's Next?</h3>
            <ul>
                <li><strong>Verify your email</strong> - Check your inbox for the verification email</li>
                <li><strong>Set your password</strong> - Complete account setup via the verification link</li>
                <li><strong>Start tracking</strong> - Log your first cannabis session in under 30 seconds</li>
                <li><strong>Discover patterns</strong> - Get insights into your consumption habits</li>
                <li><strong>Privacy first</strong> - Your data stays local unless you choose premium sync</li>
            </ul>
            
            <a href="https://cannanote.org" class="cta">Get Started with CannaNote</a>
            
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
        </div>
    </div>
</body>
</html>
		`),
		Text: fmt.Sprintf(`Welcome to CannaNote Beta

Hi there,

Welcome to the CannaNote beta! You've just joined a community focused on mindful cannabis consumption and personal wellness insights.

What's Next?
- Check your email - You'll receive a verification link shortly
- Set your password - Complete your account setup via the verification link
- Start tracking - Log your first cannabis session in under 30 seconds
- Discover patterns - Get insights into your consumption habits
- Privacy first - Your data stays local unless you choose premium sync

Get Started: https://cannanote.org

Beta Grandfathering
As a beta member, you'll receive lifetime access to premium sync features once we launch. Your early support means everything to us!

Need Help?
Visit our documentation at https://cannanote.org/docs or reply to this email.

Thank you for joining our cannabis wellness journey!
The CannaNote Team

CannaNote - Privacy-focused cannabis tracking
Privacy Policy: https://cannanote.org/privacy
Terms of Service: https://cannanote.org/terms`),
	}

	return r.sendEmail(welcomeEmail)
}

// sendEmail sends an email via Resend API
func (r *ResendService) sendEmail(email ResendEmail) error {
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp ResendResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, errorResp.Message)
		}
		return fmt.Errorf("resend API error: %d", resp.StatusCode)
	}

	return nil
}