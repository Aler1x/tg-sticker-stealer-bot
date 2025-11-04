package handlers

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tg "gopkg.in/telebot.v4"

	"tg-sticker-stiller-bot/db"
)

var adminIDs []int64

func InitAdminIDs() {
	adminIDsStr := os.Getenv("ADMIN_IDS")
	if adminIDsStr == "" {
		log.Println("Warning: ADMIN_IDS not set, broadcast feature will be disabled")
		return
	}

	ids := strings.Split(adminIDsStr, ",")
	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Warning: Invalid admin ID '%s': %v", idStr, err)
			continue
		}
		adminIDs = append(adminIDs, id)
	}

	if len(adminIDs) > 0 {
		log.Printf("Loaded %d admin ID(s)", len(adminIDs))
	}
}

func IsAdmin(userID int64) bool {
	for _, id := range adminIDs {
		if id == userID {
			return true
		}
	}
	return false
}

func HandleAdminStats(ctx tg.Context, repo *db.Repository) error {
	if !IsAdmin(ctx.Sender().ID) {
		return nil
	}

	packCount, err := repo.GetPackCount()
	if err != nil {
		log.Printf("Failed to get pack count: %v", err)
		return ctx.Send("❌ Failed to fetch statistics.")
	}

	userCount, err := repo.GetUserCount()
	if err != nil {
		log.Printf("Failed to get user count: %v", err)
		return ctx.Send("❌ Failed to fetch statistics.")
	}

	stats := fmt.Sprintf(
		"*Bot Statistics*\n\n"+
			"Active users: `%d`\n"+
			"Total packs: `%d`",
		userCount,
		packCount,
	)

	return ctx.Send(stats, &tg.SendOptions{ParseMode: tg.ModeMarkdown})
}
