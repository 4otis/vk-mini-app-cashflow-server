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

func (s *GameService) PlayerIsReady(ctx context.Context, code string, VKID int) ([]dto.PlayerResponse, error) {
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
		return nil, err
	}

	// if !player.Ready {

	// 	for _, pl := range players {
	// 		if !pl.Ready {
	// 			break
	// 		}
	// 	}

	// 	log.Printf("ALL PLAYERS ARE READY!!!")
	// 	log.Printf("ALL PLAYERS ARE READY!!!")
	// 	log.Printf("ALL PLAYERS ARE READY!!!")

	// }

	result := make([]dto.PlayerResponse, 0, len(players))
	for _, p := range players {
		result = append(result, *convertPlayerToDTO(&p))
	}

	return result, nil
}

func (s *GameService) ArePlayersReady(ctx context.Context, code string) (dto.PlayersAreReadyRespone, error) {
	var response dto.PlayersAreReadyRespone
	response.Ready = false

	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return response, fmt.Errorf("failed to get session: %w", err)
	}

	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return response, fmt.Errorf("failed to read all players: %w", err)
	}

	for _, pl := range players {
		if !pl.Ready {
			return response, nil
		}
	}

	log.Printf("ALL PLAYERS ARE READY!!!")
	log.Printf("ALL PLAYERS ARE READY!!!")
	log.Printf("ALL PLAYERS ARE READY!!!")

	response.Ready = true

	return response, nil
}

func (s *GameService) InitGameState(ctx context.Context, code string) (dto.GameStateResponse, error) {
	var response dto.GameStateResponse

	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return response, err
	}

	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return response, err
	}

	response.SessionCode = code
	response.CurrentTurn = 0
	response.Players = []dto.PlayerStat{}

	for _, p := range players {
		err = s.playerRepo.InitPlayer(p.ID)
		if err != nil {
			return response, err
		}

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
			BankLoan:      p.BankLoan,
		})
	}

	return response, nil
}

func (s *GameService) LoadGameState(ctx context.Context, code string) (dto.GameStateResponse, error) {
	var response dto.GameStateResponse

	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return response, err
	}

	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return response, err
	}

	response.SessionCode = code
	response.CurrentTurn = 0
	response.Players = []dto.PlayerStat{}

	for _, p := range players {
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
			BankLoan:      p.BankLoan,
		})
	}

	return response, nil
}

func (s *GameService) EndTurn(ctx context.Context, code string, VKID int) error {
	session, err := s.sessionRepo.Read(code)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	players, err := s.playerRepo.ReadAll(session.ID)
	if err != nil {
		return err
	}

	cnt := len(players)
	newTurn := int(session.CurrentTurn+1) % cnt

	err = s.sessionRepo.UpdateFields(session.ID, map[string]interface{}{
		"current_turn": newTurn,
	})
	if err != nil {
		return err
	}
	return nil
}
