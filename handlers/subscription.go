package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/utils"

	tele "gopkg.in/telebot.v4"
)

// ShowSubscriptionPlans displays all available subscription options
func ShowSubscriptionPlans(c tele.Context, repo *db.Repository) error {
	prices, err := repo.GetAllSubscriptionPrices()
	if err != nil {
		utils.LogError("Failed to get subscription prices", err)
		return c.Send("Failed to load subscription plans. Please try again later.")
	}

	// Build message with all plans
	var message strings.Builder
	message.WriteString("💎 *Emoji Stealing Subscription Plans*\n\n")
	message.WriteString("Choose a plan to start stealing emoji:\n\n")

	// Create inline keyboard with subscription options
	var buttons [][]tele.InlineButton

	for _, price := range prices {
		emoji := getSubscriptionEmoji(price.SubscriptionType)
		description := fmt.Sprintf("%s %s - %d ⭐", emoji, price.Description, price.PriceStars)
		message.WriteString(fmt.Sprintf("• %s\n", description))

		btn := tele.InlineButton{
			Text: description,
			Data: fmt.Sprintf("sub_%s", price.SubscriptionType),
		}
		buttons = append(buttons, []tele.InlineButton{btn})
	}

	message.WriteString("\n_Prices in Telegram Stars ⭐_")

	markup := &tele.ReplyMarkup{
		InlineKeyboard: buttons,
	}

	return c.Send(message.String(), markup, tele.ModeMarkdown)
}

// HandleSubscriptionCallback handles subscription button clicks
func HandleSubscriptionCallback(c tele.Context, repo *db.Repository) error {
	data := c.Callback().Data

	if !strings.HasPrefix(data, "sub_") {
		return nil
	}

	subType := db.SubscriptionType(strings.TrimPrefix(data, "sub_"))

	// Get price info
	price, err := repo.GetSubscriptionPrice(subType)
	if err != nil {
		utils.LogError("Failed to get subscription price", err)
		return c.Respond(&tele.CallbackResponse{Text: "Error loading price"})
	}
	if price == nil {
		return c.Respond(&tele.CallbackResponse{Text: "Invalid subscription type"})
	}

	// Create invoice
	invoice := &tele.Invoice{
		Title:       price.Description,
		Description: fmt.Sprintf("Subscribe to steal emoji with %s plan", price.Description),
		Payload:     fmt.Sprintf("subscription_%s", subType),
		Currency:    "XTR", // Telegram Stars
		Prices: []tele.Price{
			{Label: price.Description, Amount: price.PriceStars},
		},
	}

	// Send invoice
	_, err = c.Bot().Send(c.Sender(), invoice)
	if err != nil {
		utils.LogError("Failed to send invoice", err)
		return c.Respond(&tele.CallbackResponse{Text: "Failed to create payment"})
	}

	return c.Respond(&tele.CallbackResponse{Text: "Payment invoice sent!"})
}

// HandlePreCheckoutQuery handles pre-checkout validation
func HandlePreCheckoutQuery(c tele.Context, repo *db.Repository) error {
	query := c.PreCheckoutQuery()

	// Validate payload
	if !strings.HasPrefix(query.Payload, "subscription_") {
		return c.Bot().Accept(query, "Invalid subscription")
	}

	// Accept the checkout
	return c.Bot().Accept(query)
}

// HandleSuccessfulPayment handles successful payment
func HandleSuccessfulPayment(c tele.Context, repo *db.Repository) error {
	payment := c.Message().Payment

	// Parse subscription type from payload
	if !strings.HasPrefix(payment.Payload, "subscription_") {
		return c.Send("Invalid payment payload")
	}

	subTypeStr := strings.TrimPrefix(payment.Payload, "subscription_")
	subType := db.SubscriptionType(subTypeStr)

	// Get price info to determine subscription details
	price, err := repo.GetSubscriptionPrice(subType)
	if err != nil {
		utils.LogError("Failed to get subscription price", err)
		return c.Send("Payment received but failed to activate subscription. Please contact support.")
	}
	if price == nil {
		return c.Send("Invalid subscription type")
	}

	// Create subscription based on type
	userSub := &db.UserSubscription{
		UserID:           c.Sender().ID,
		SubscriptionType: subType,
	}

	// Determine if it's count-based, time-based, or infinity
	switch subType {
	case db.SubscriptionOneSteal, db.SubscriptionTenSteals:
		// Count-based subscription
		count := price.Value
		userSub.RemainingCount = &count
		userSub.ExpiresAt = nil
	case db.SubscriptionWeek, db.SubscriptionMonth, db.SubscriptionYear:
		// Time-based subscription
		expiresAt := time.Now().AddDate(0, 0, price.Value)
		userSub.ExpiresAt = &expiresAt
		userSub.RemainingCount = nil
	case db.SubscriptionInfinity:
		// Infinity subscription - set expiry to 100 years from now
		expiresAt := time.Now().AddDate(100, 0, 0)
		userSub.ExpiresAt = &expiresAt
		userSub.RemainingCount = nil
	}

	// Save subscription
	err = repo.CreateUserSubscription(userSub)
	if err != nil {
		utils.LogError("Failed to create user subscription", err)
		return c.Send("Payment received but failed to activate subscription. Please contact support.")
	}

	// Record payment history
	chargeID := payment.ChargeID
	provider := payment.ProviderChargeID
	paymentHistory := &db.PaymentHistory{
		UserID:            c.Sender().ID,
		SubscriptionType:  subType,
		PriceStars:        payment.Total,
		PaymentChargeID:   &chargeID,
		PaymentProvider:   &provider,
		Status:            "completed",
	}

	err = repo.CreatePaymentHistory(paymentHistory)
	if err != nil {
		utils.LogError("Failed to record payment history", err)
		// Don't fail the whole process if history recording fails
	}

	// Send confirmation
	var confirmMessage string
	if userSub.RemainingCount != nil {
		confirmMessage = fmt.Sprintf(
			"✅ *Payment Successful!*\n\n"+
				"Your subscription is now active!\n"+
				"Remaining steals: %d\n\n"+
				"You can now steal emoji by sending me an emoji pack link!",
			*userSub.RemainingCount,
		)
	} else if subType == db.SubscriptionInfinity {
		confirmMessage = "✅ *Payment Successful!*\n\n" +
			"Your *LIFETIME* subscription is now active!\n" +
			"♾️ Unlimited emoji steals forever!\n\n" +
			"You can now steal emoji by sending me an emoji pack link!"
	} else {
		confirmMessage = fmt.Sprintf(
			"✅ *Payment Successful!*\n\n"+
				"Your subscription is now active!\n"+
				"Valid until: %s\n\n"+
				"You can now steal emoji by sending me an emoji pack link!",
			userSub.ExpiresAt.Format("2006-01-02 15:04"),
		)
	}

	return c.Send(confirmMessage, tele.ModeMarkdown)
}

// CheckSubscription checks if user has active subscription
func CheckSubscription(userID int64, repo *db.Repository) (*db.UserSubscription, error) {
	sub, err := repo.GetActiveSubscription(userID)
	if err != nil {
		return nil, err
	}

	if sub == nil {
		return nil, nil
	}

	// Double check validity
	if sub.RemainingCount != nil && *sub.RemainingCount <= 0 {
		return nil, nil
	}

	if sub.ExpiresAt != nil && sub.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return sub, nil
}

// ConsumeSubscription decrements subscription count or checks expiry
func ConsumeSubscription(sub *db.UserSubscription, repo *db.Repository) error {
	// For count-based subscriptions, decrement the count
	if sub.RemainingCount != nil {
		return repo.DecrementSubscriptionCount(sub.ID)
	}

	// For time-based subscriptions, just check expiry (already checked in CheckSubscription)
	return nil
}

func getSubscriptionEmoji(subType db.SubscriptionType) string {
	switch subType {
	case db.SubscriptionOneSteal:
		return "🎯"
	case db.SubscriptionTenSteals:
		return "🎁"
	case db.SubscriptionWeek:
		return "📅"
	case db.SubscriptionMonth:
		return "📆"
	case db.SubscriptionYear:
		return "🎊"
	case db.SubscriptionInfinity:
		return "♾️"
	default:
		return "💎"
	}
}

// HandleCheckMySubscription shows user's current subscription status
func HandleCheckMySubscription(c tele.Context, repo *db.Repository) error {
	sub, err := CheckSubscription(c.Sender().ID, repo)
	if err != nil {
		utils.LogError("Failed to check subscription", err)
		return c.Send("Failed to check subscription status")
	}

	if sub == nil {
		return c.Send(
			"❌ You don't have an active subscription.\n\n"+
				"Use /subscribe to view available plans!",
		)
	}

	var message string
	emoji := getSubscriptionEmoji(sub.SubscriptionType)

	if sub.RemainingCount != nil {
		message = fmt.Sprintf(
			"%s *Active Subscription*\n\n"+
				"Plan: %s\n"+
				"Remaining steals: %d\n"+
				"Activated: %s",
			emoji,
			sub.SubscriptionType,
			*sub.RemainingCount,
			sub.CreatedAt.Format("2006-01-02 15:04"),
		)
	} else if sub.SubscriptionType == db.SubscriptionInfinity {
		message = fmt.Sprintf(
			"%s *LIFETIME Subscription*\n\n"+
				"Plan: %s\n"+
				"Status: ♾️ Unlimited Forever\n"+
				"Activated: %s",
			emoji,
			sub.SubscriptionType,
			sub.CreatedAt.Format("2006-01-02 15:04"),
		)
	} else {
		message = fmt.Sprintf(
			"%s *Active Subscription*\n\n"+
				"Plan: %s\n"+
				"Valid until: %s\n"+
				"Activated: %s",
			emoji,
			sub.SubscriptionType,
			sub.ExpiresAt.Format("2006-01-02 15:04"),
			sub.CreatedAt.Format("2006-01-02 15:04"),
		)
	}

	return c.Send(message, tele.ModeMarkdown)
}

// Admin command to grant subscription manually
func HandleAdminGrantSubscription(c tele.Context, repo *db.Repository) error {
	if !IsAdmin(c.Sender().ID) {
		return c.Send("⛔ This command is only available to administrators")
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Send(
			"Usage: /grant <user_id> <subscription_type> [value]\n\n"+
				"Subscription types:\n"+
				"• one_steal, ten_steals - count based\n"+
				"• week, month, year - time based\n"+
				"• infinity - lifetime unlimited\n\n"+
				"Examples:\n"+
				"/grant 123456789 month\n"+
				"/grant 123456789 ten_steals\n"+
				"/grant 123456789 one_steal 5 (custom count)\n"+
				"/grant 123456789 infinity",
		)
	}

	// Parse user ID
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return c.Send("Invalid user ID")
	}

	// Parse subscription type
	subType := db.SubscriptionType(args[1])

	// Get price info for default values
	price, err := repo.GetSubscriptionPrice(subType)
	if err != nil || price == nil {
		return c.Send("Invalid subscription type")
	}

	// Create subscription
	userSub := &db.UserSubscription{
		UserID:           userID,
		SubscriptionType: subType,
	}

	// Determine subscription details
	switch subType {
	case db.SubscriptionOneSteal, db.SubscriptionTenSteals:
		// Count-based
		count := price.Value
		// Allow custom count from args
		if len(args) >= 3 {
			customCount, err := strconv.Atoi(args[2])
			if err == nil && customCount > 0 {
				count = customCount
			}
		}
		userSub.RemainingCount = &count
	case db.SubscriptionWeek, db.SubscriptionMonth, db.SubscriptionYear:
		// Time-based
		expiresAt := time.Now().AddDate(0, 0, price.Value)
		userSub.ExpiresAt = &expiresAt
	case db.SubscriptionInfinity:
		// Infinity subscription - set expiry to 100 years from now
		expiresAt := time.Now().AddDate(100, 0, 0)
		userSub.ExpiresAt = &expiresAt
	default:
		return c.Send("Invalid subscription type")
	}

	// Save subscription
	err = repo.CreateUserSubscription(userSub)
	if err != nil {
		utils.LogError("Failed to grant subscription", err)
		return c.Send("Failed to grant subscription")
	}

	// Record in payment history
	paymentHistory := &db.PaymentHistory{
		UserID:           userID,
		SubscriptionType: subType,
		PriceStars:       0, // Admin grant = free
		Status:           "admin_granted",
	}
	repo.CreatePaymentHistory(paymentHistory)

	var confirmMsg string
	if userSub.RemainingCount != nil {
		confirmMsg = fmt.Sprintf(
			"✅ Subscription granted!\n\n"+
				"User ID: %d\n"+
				"Type: %s\n"+
				"Remaining: %d steals",
			userID, subType, *userSub.RemainingCount,
		)
	} else if subType == db.SubscriptionInfinity {
		confirmMsg = fmt.Sprintf(
			"✅ LIFETIME Subscription granted!\n\n"+
				"User ID: %d\n"+
				"Type: %s\n"+
				"Status: ♾️ Unlimited Forever",
			userID, subType,
		)
	} else {
		confirmMsg = fmt.Sprintf(
			"✅ Subscription granted!\n\n"+
				"User ID: %d\n"+
				"Type: %s\n"+
				"Expires: %s",
			userID, subType, userSub.ExpiresAt.Format("2006-01-02 15:04"),
		)
	}

	return c.Send(confirmMsg)
}

// Admin command to update subscription prices
func HandleAdminSetPrice(c tele.Context, repo *db.Repository) error {
	if !IsAdmin(c.Sender().ID) {
		return c.Send("⛔ This command is only available to administrators")
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Send(
			"Usage: /setprice <subscription_type> <price_stars>\n\n"+
				"Subscription types:\n"+
				"• one_steal\n"+
				"• ten_steals\n"+
				"• week\n"+
				"• month\n"+
				"• year\n"+
				"• infinity\n\n"+
				"Example: /setprice month 200",
		)
	}

	subType := db.SubscriptionType(args[0])
	priceStars, err := strconv.Atoi(args[1])
	if err != nil || priceStars < 0 {
		return c.Send("Invalid price. Must be a positive number.")
	}

	err = repo.UpdateSubscriptionPrice(subType, priceStars)
	if err != nil {
		utils.LogError("Failed to update price", err)
		return c.Send("Failed to update price. Make sure subscription type is valid.")
	}

	return c.Send(fmt.Sprintf("✅ Price updated!\n%s: %d ⭐", subType, priceStars))
}

// HandleAdminViewPrices shows all current prices
func HandleAdminViewPrices(c tele.Context, repo *db.Repository) error {
	if !IsAdmin(c.Sender().ID) {
		return c.Send("⛔ This command is only available to administrators")
	}

	prices, err := repo.GetAllSubscriptionPrices()
	if err != nil {
		utils.LogError("Failed to get prices", err)
		return c.Send("Failed to load prices")
	}

	var message strings.Builder
	message.WriteString("💰 *Current Subscription Prices*\n\n")

	for _, price := range prices {
		emoji := getSubscriptionEmoji(price.SubscriptionType)
		message.WriteString(fmt.Sprintf(
			"%s %s: *%d ⭐*\n   (%s)\n\n",
			emoji, price.SubscriptionType, price.PriceStars, price.Description,
		))
	}

	message.WriteString("Use /setprice to update prices")

	return c.Send(message.String(), tele.ModeMarkdown)
}
