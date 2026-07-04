package services

import (
	"sync"
	"tg-sticker-stiller-bot/types"

	tg "gopkg.in/telebot.v4"
)

type SessionState string

const (
	StateIdle               SessionState = ""
	StateWaitingForPackName SessionState = "waiting_for_pack_name"
)

type PackAction string

const (
	ActionCopy      PackAction = "copy"
	ActionDownload  PackAction = "download"
	ActionTransform PackAction = "transform"
)

type Session struct {
	State            SessionState
	Action           PackAction
	OriginalPackName string
	OriginalItems    []tg.Sticker
	Title            string
	Name             string
	FullLink         string
	PackType         types.StickerType
	ProgressMsgID    int
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[int64]*Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[int64]*Session),
	}
}

func (s *SessionStore) Get(userID int64) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if session, exists := s.sessions[userID]; exists {
		return session
	}

	return &Session{State: StateIdle}
}

func (s *SessionStore) Set(userID int64, session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[userID] = session
}

func (s *SessionStore) Clear(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, userID)
}
