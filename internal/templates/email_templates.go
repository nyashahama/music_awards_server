// Package templates
package templates

const (
	WelcomeEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4F46E5; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 8px 8px; }
        .button { background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Music Awards!</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Username}},</h2>
            <p>Welcome to the Music Awards platform! Your account has been successfully created.</p>
            <p>You can now:</p>
            <ul>
                <li>Vote for your favorite artists in various categories</li>
                <li>Track your voting history</li>
                <li>See real-time results</li>
            </ul>
            <p>You have <strong>{{.VoteCount}} votes</strong> available to use.</p>
            <p>
                <a href="{{.AppURL}}" class="button">Start Voting Now</a>
            </p>
            <p>If you have any questions, please contact our support team.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.CurrentYear}} Music Awards. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	PasswordResetTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #DC2626; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 8px 8px; }
        .button { background: #DC2626; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
        .token { background: #f0f0f0; padding: 10px; border-radius: 4px; font-family: monospace; margin: 10px 0; }
        .warning { color: #DC2626; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Hello,</h2>
            <p>We received a request to reset your password for your Music Awards account.</p>
            <p>Click the button below to reset your password:</p>
            <p>
                <a href="{{.ResetURL}}" class="button">Reset Password</a>
            </p>
            <p>Or copy and paste this token in the reset form:</p>
            <div class="token">{{.Token}}</div>
            <p class="warning">This token will expire in 1 hour for security reasons.</p>
            <p>If you didn't request this reset, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.CurrentYear}} Music Awards. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	LoginNotificationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #059669; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 20px; border-radius: 0 0 8px 8px; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
        .info-box { background: #e8f5e8; padding: 15px; border-radius: 4px; margin: 15px 0; }
        .warning { color: #DC2626; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>New Login Detected</h1>
        </div>
        <div class="content">
            <h2>Hello {{.Username}},</h2>
            <p>We noticed a recent login to your Music Awards account:</p>
            <div class="info-box">
                <p><strong>Time:</strong> {{.LoginTime}}</p>
                <p><strong>Device/Browser:</strong> {{.UserAgent}}</p>
                <p><strong>IP Address:</strong> {{.IPAddress}}</p>
            </div>
            <p>If this was you, you can safely ignore this email.</p>
            <p class="warning">If you don't recognize this activity, please reset your password immediately and contact our support team.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.CurrentYear}} Music Awards. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
)
