# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Migrated from SQLite to PostgreSQL for better scalability and production deployment
- Updated database driver from `gorm.io/driver/sqlite` to `gorm.io/driver/postgres`
- Changed `DB_PATH` environment variable to `DATABASE_URL` for PostgreSQL connection string
- Removed SQLite-specific dependencies from Dockerfile
- Updated documentation to reflect PostgreSQL setup

### Removed
- Pack tracking system (`/list` and `/delete` commands)
- Pack storage in database (packs are only created on Telegram, not stored in bot database)
- `Pack` model and `PackRepository`

## [1.0.0] - 2025-10-26

### Initial realease of bot

## Features

- Copy sticker packs
- Copy emoji packs
- Download packs as ZIP
- Convert images to stickers
- Convert stickers to images
- Settings for default action
- Multi-language support (English, Ukrainian, Polish)
- View bot statistics (admin only)

## Internal changes

- Use webhooks for production (Railway)
- Use long polling for local development
- Use SQLite database for persistence
- Use Telebotv4 framework for Telegram API
- Use Docker for deployment
