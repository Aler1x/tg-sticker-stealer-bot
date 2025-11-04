package handlers

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	tg "gopkg.in/telebot.v4"

	"tg-sticker-stiller-bot/db"
	"tg-sticker-stiller-bot/utils"
)

var adminIDs []int64

func InitAdminIDs() {
	adminIDsStr := os.Getenv("ADMIN_IDS")
	if adminIDsStr == "" {
		utils.Logger("warn", "ADMIN_IDS not set, broadcast feature will be disabled")
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
			utils.Logger("warn", "Invalid admin ID", map[string]any{
				"adminId": idStr,
				"error":   err.Error(),
			})
			continue
		}
		adminIDs = append(adminIDs, id)
	}

	if len(adminIDs) > 0 {
		utils.Logger("info", "Admin IDs loaded", map[string]any{"count": len(adminIDs)})
	}
}

func IsAdmin(userID int64) bool {
	return slices.Contains(adminIDs, userID)
}

func HandleAdminStats(ctx tg.Context, repo *db.Repository) error {
	if !IsAdmin(ctx.Sender().ID) {
		return nil
	}

	packCount, err := repo.GetPackCount()
	if err != nil {
		utils.Logger("error", "Failed to get pack count", map[string]any{"error": err.Error()})
		return ctx.Send("❌ Failed to fetch statistics.")
	}

	userCount, err := repo.GetUserCount()
	if err != nil {
		utils.Logger("error", "Failed to get user count", map[string]any{"error": err.Error()})
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
