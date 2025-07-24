package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/dto"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/models"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidCode     = errors.New("invalid session code")
)

type SessionService struct {
	sessionRepo repository.SessionRepository
	playerRepo  repository.PlayerRepository
}

func NewSessionService(sessionRepo *repository.SessionRepository, playerRepo *repository.PlayerRepository) *SessionService {
	return &SessionService{
		sessionRepo: *sessionRepo,
		playerRepo:  *playerRepo,
	}
}

// Generate random session code
func generateSessionCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	const codeLength = 6
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *SessionService) CreateSession(ctx context.Context, creatorVKID int, nickname string) (*dto.CreateSessionResponse, error) {
	// Create session
	session := &models.Session{
		Code:     generateSessionCode(),
		IsActive: false,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create creator player
	player := &models.Player{
		VKID:      creatorVKID,
		SessionID: session.ID,
		Nickname:  nickname,
		Ready:     true,
		Balance:   10000, // Starting balance
	}

	if err := s.playerRepo.Create(player); err != nil {
		// Rollback session creation if player creation fails
		_ = s.sessionRepo.Delete(session.ID)
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return &dto.CreateSessionResponse{
		Code:     session.Code,
		JoinLink: fmt.Sprintf("/join/%s", strconv.Itoa((int(session.ID)))),
	}, nil
}

func (s *SessionService) JoinSession(ctx context.Context, code string, vkID int, nickname string) (*dto.PlayerResponse, error) {
	session, err := s.sessionRepo.Read(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if player already exists in this session
	existingPlayers, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing players: %w", err)
	}

	for _, p := range existingPlayers {
		if p.VKID == vkID {
			return convertPlayerToDTO(&p), nil
		}
	}

	// Create new player
	player := &models.Player{
		VKID:      vkID,
		SessionID: session.ID,
		Nickname:  nickname,
		Ready:     false,
		Balance:   10000, // Starting balance
	}

	if err := s.playerRepo.Create(player); err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return convertPlayerToDTO(player), nil
}

func (s *SessionService) GetSessionPlayers(ctx context.Context, code string) ([]dto.PlayerResponse, error) {
	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	result := make([]dto.PlayerResponse, 0, len(players))
	for _, p := range players {
		result = append(result, *convertPlayerToDTO(&p))
	}

	return result, nil
}

func convertPlayerToDTO(player *models.Player) *dto.PlayerResponse {
	return &dto.PlayerResponse{
		ID:            player.ID,
		VKID:          player.VKID,
		Nickname:      player.Nickname,
		Ready:         player.Ready,
		Position:      player.Position,
		CharacterID:   player.CharacterID,
		PassiveIncome: player.PassiveIncome,
		TotalIncome:   player.TotalIncome,
		CashFlow:      player.CashFlow,
		Balance:       player.Balance,
		BankLoan:      player.BankLoan,
	}
}
