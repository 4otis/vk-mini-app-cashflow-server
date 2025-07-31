package services

import (
	"context"
	"fmt"
	"log"

	"github.com/4otis/vk-mini-app-cashflow-server/internal/dto"
	"github.com/4otis/vk-mini-app-cashflow-server/internal/repository"
)

type GameService struct {
	sessionRepo repository.SessionRepository
	playerRepo  repository.PlayerRepository
	assetRepo   repository.AssetRepository
}

func NewGameService(sessionRepo *repository.SessionRepository,
	playerRepo *repository.PlayerRepository,
	assetRepo *repository.AssetRepository) *GameService {
	return &GameService{
		sessionRepo: *sessionRepo,
		playerRepo:  *playerRepo,
		assetRepo:   *assetRepo,
	}
}

func (s *GameService) TryStartGame(ctx context.Context, code string, VKID int) ([]dto.PlayerResponse, error) {
	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	player, err := s.playerRepo.ReadByVKID(VKID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	err = s.playerRepo.UpdateFields(player.ID, map[string]interface{}{
		"ready": !player.Ready,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}
	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {

	}

	if !player.Ready {

		for _, pl := range players {
			if !pl.Ready {
				break
			}
		}

		log.Printf("ALL PLAYERS ARE READY!!!")
		log.Printf("ALL PLAYERS ARE READY!!!")
		log.Printf("ALL PLAYERS ARE READY!!!")

	}

	result := make([]dto.PlayerResponse, 0, len(players))
	for _, p := range players {
		result = append(result, *convertPlayerToDTO(&p))
	}

	return result, nil
}

func (s *GameService) InitPlayers(ctx context.Context, code string) (dto.GameStateResponse, error) {
	var response dto.GameStateResponse

	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return response, err
	}

	log.Println("babah")

	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return response, err
	}
	log.Println("babah2")

	response.SessionCode = code
	response.CurrentTurn = 0
	response.Players = []dto.PlayerStat{}

	for _, p := range players {
		err = s.playerRepo.InitPlayer(p.ID)
		if err != nil {
			return response, err
		}
		log.Println("babah3")

		response.Players = append(response.Players, dto.PlayerStat{
			VKID:          p.VKID,
			Nickname:      p.Nickname,
			PassiveIncome: p.PassiveIncome,
			TotalIncome:   p.TotalIncome,
			TotalExpenses: p.TotalExpenses,
			Cashflow:      p.Cashflow,
			Position:      p.Position,
			Balance:       p.Balance,
			ChildAmount:   p.ChildAmount,
		})
	}

	return response, nil
}
