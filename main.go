package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/handlers"
	"tg-sticker-stiller-bot/services"
	"tg-sticker-stiller-bot/types"
	"tg-sticker-stiller-bot/utils"
	"time"

	tg "gopkg.in/telebot.v4"
)

func main() {
	log.Println("Starting bot...")

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN environment variable is not set")
	}

	if err := utils.EnsureTempDir(); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	repo, err := db.NewRepository("./data/packs.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	// Configure poller based on environment
	var poller tg.Poller
	publicURL := os.Getenv("PUBLIC_URL")

	if publicURL != "" {
		// Use webhooks for production (Railway)
		port := os.Getenv("PORT")
		if port == "" {
			port = "8443"
		}

		webhookURL := publicURL + "/webhook"
		log.Printf("Using webhook mode: %s", webhookURL)

		poller = &tg.Webhook{
			Listen:   "0.0.0.0:" + port,
			Endpoint: &tg.WebhookEndpoint{PublicURL: webhookURL},
		}
	} else {
		// Use long polling for local development
		log.Println("Using long polling mode (local development)")
		poller = &tg.LongPoller{Timeout: 10 * time.Second}
	}

	bot, err := tg.NewBot(tg.Settings{
		Token:  token,
		Poller: poller,
	})
	utils.FailFast(err)

	name := bot.Me.Username
	sessions := services.NewSessionStore()

	handlers.InitAdminIDs()

	bot.Use(tg.MiddlewareFunc(func(next tg.HandlerFunc) tg.HandlerFunc {
		return func(ctx tg.Context) error {
			if ctx.Message() != nil {
				log.Printf("User: %d, Message ID: %d, Text: %s",
					ctx.Sender().ID, ctx.Message().ID, ctx.Message().Text)
			}
			return next(ctx)
		}
	}))

	bot.Use(tg.MiddlewareFunc(func(next tg.HandlerFunc) tg.HandlerFunc {
		return func(ctx tg.Context) error {
			if ctx.Sender() != nil {
				user := &db.User{
					UserID:       ctx.Sender().ID,
					Username:     ctx.Sender().Username,
					FirstName:    ctx.Sender().FirstName,
					LastName:     ctx.Sender().LastName,
					LanguageCode: ctx.Sender().LanguageCode,
				}
				if err := repo.UpsertUser(user); err != nil {
					log.Printf("Failed to track user %d: %v", ctx.Sender().ID, err)
				}
			}
			return next(ctx)
		}
	}))

	bot.SetCommands([]tg.Command{
		{Text: "/start", Description: utils.T("en", "start-command")},
		{Text: "/help", Description: utils.T("en", "help-command")},
		{Text: "/list", Description: utils.T("en", "list-command")},
		{Text: "/delete", Description: utils.T("en", "delete-command")},
		{Text: "/cancel", Description: "Cancel current operation"},
	})

	bot.Handle("/start", func(ctx tg.Context) error {
		lang := ctx.Message().Sender.LanguageCode
		username := ctx.Message().Sender.Username
		sessions.Clear(ctx.Sender().ID)
		return ctx.Send(utils.T(lang, "welcome", username))
	})

	bot.Handle("/help", func(ctx tg.Context) error {
		lang := ctx.Message().Sender.LanguageCode
		return ctx.Send(utils.T(lang, "help"))
	})

	bot.Handle("/list", func(ctx tg.Context) error {
		return handlers.HandleListPacks(ctx, repo)
	})

	bot.Handle("/delete", func(ctx tg.Context) error {
		lang := ctx.Message().Sender.LanguageCode
		args := strings.Fields(ctx.Text())
		if len(args) < 2 {
			return ctx.Send(utils.T(lang, "delete-usage"))
		}

		packID, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return ctx.Send(utils.T(lang, "delete-usage"))
		}

		return handlers.HandleDeletePack(ctx, packID, repo)
	})

	bot.Handle("/cancel", func(ctx tg.Context) error {
		lang := ctx.Message().Sender.LanguageCode
		userID := ctx.Sender().ID
		session := sessions.Get(userID)

		if session.State == services.StateIdle {
			return ctx.Send(utils.T(lang, "help"))
		}

		sessions.Clear(userID)
		return ctx.Send(utils.T(lang, "cancelled"))
	})

	bot.Handle("/stats", func(ctx tg.Context) error {
		return handlers.HandleAdminStats(ctx, repo)
	})

	bot.Handle(tg.OnText, func(ctx tg.Context) error {
		text := ctx.Text()
		userID := ctx.Sender().ID
		lang := ctx.Message().Sender.LanguageCode

		session := sessions.Get(userID)

		switch session.State {
		case services.StateWaitingForPackName:
			return handlers.HandlePackNameInput(ctx, text, bot, sessions, repo)

		default:
			if utils.IsStickerPack(text) {
				packName := utils.ExtractStickerPackName(text)
				if packName == "" {
					return ctx.Send(utils.T(lang, "invalid-link"))
				}
				return handlers.HandlePack(ctx, packName, types.StickerTypeRegular, bot, sessions)
			}

			if utils.IsEmojiPack(text) {
				packName := utils.ExtractEmojiPackName(text)
				if packName == "" {
					return ctx.Send(utils.T(lang, "invalid-link"))
				}
				return handlers.HandlePack(ctx, packName, types.StickerTypeEmoji, bot, sessions)
			}

			return ctx.Send(utils.T(lang, "invalid-link"))
		}
	})

	go func() {
		log.Printf("Bot @%s started successfully\n", name)
		if publicURL != "" {
			log.Printf("Webhook endpoint: %s/webhook", publicURL)
		}
		bot.Start()
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println("Received an interrupt, stopping...")

	bot.Stop()
	log.Println("Bot stopped")
}
