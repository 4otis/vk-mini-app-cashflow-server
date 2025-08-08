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
	marketRepo  repository.MarketRepository
	issueRepo   repository.IssueRepository
}

func NewGameService(sessionRepo *repository.SessionRepository,
	playerRepo *repository.PlayerRepository,
	assetRepo *repository.AssetRepository,
	marketRepo *repository.MarketRepository,
	issueRepo *repository.IssueRepository) *GameService {
	return &GameService{
		sessionRepo: *sessionRepo,
		playerRepo:  *playerRepo,
		assetRepo:   *assetRepo,
		marketRepo:  *marketRepo,
		issueRepo:   *issueRepo,
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
	response.CurrentTurn = session.CurrentTurn
	response.Players = []dto.PlayerStat{}

	for _, p := range players {
		// TODO: получаю assets конкретного игрока по p.ID
		assetsDb, err := s.assetRepo.ReadAllByPlayerID(p.ID)
		if err != nil {
			return response, err
		}

		var assets []dto.AssetStat
		for _, a := range assetsDb {
			assets = append(assets, dto.AssetStat{
				Title:    a.Title,
				Price:    a.Price,
				Cashflow: a.Cashflow,
			})
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
			// TODO: добавляю сформированный список []AssetStat
			Assets: assets,
		})
	}

	return response, nil
}

func chooseCardType() string {
	switch repository.RandRange(0, 2) {
	case 0:
		return "ASSET"
	case 1:
		return "MARKET"
	case 2:
		return "ISSUE"
	case 3:
		return "CHILD"
	}

	return ""
}

func (s *GameService) RollDice(ctx context.Context, code string, VKID int, value int) (dto.RollDiceResponse, error) {
	var resp dto.RollDiceResponse

	player, err := s.playerRepo.MovePlayer(VKID, value)
	if err != nil {
		return resp, err
	}

	resp.Player = dto.PlayerStat{
		VKID:          player.VKID,
		Nickname:      player.Nickname,
		PassiveIncome: player.PassiveIncome,
		TotalIncome:   player.TotalIncome,
		TotalExpenses: player.TotalExpenses,
		Cashflow:      player.Cashflow,
		Position:      player.Position,
	}

	CardType := chooseCardType()
	switch CardType {
	case "ASSET":
		card, err := s.assetRepo.ReadRandom()
		if err != nil {
			return resp, err
		}

		resp.CurrentCard = dto.Card{
			Type: CardType,
			Asset: dto.AssetCard{
				Title:    card.Title,
				Descr:    card.Descr,
				Price:    card.Price,
				Cashflow: card.Cashflow,
			},
		}
	case "MARKET":
		card, err := s.marketRepo.ReadRandom()
		if err != nil {
			return resp, err
		}

		resp.CurrentCard = dto.Card{
			Type: CardType,
			Market: dto.MarketCard{
				Title:    card.Title,
				Descr:    card.Descr,
				SellCost: card.SellCost,
			},
		}
	case "ISSUE":
		card, err := s.issueRepo.ReadRandom()
		if err != nil {
			return resp, err
		}

		resp.CurrentCard = dto.Card{
			Type: CardType,
			Issue: dto.IssueCard{
				Title: card.Title,
				Descr: card.Descr,
				Price: card.Price,
			},
		}
	default:
		return resp, fmt.Errorf("failed to indentify card")
	}

	return resp, nil
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

func (s *GameService) BuyAsset(ctx context.Context, code string, req *dto.CardActionBuyReq) error {
	err := s.playerRepo.BuyAsset(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *GameService) SellAsset(ctx context.Context, code string, req *dto.CardActionSellReq) error {
	err := s.playerRepo.SellAsset(req)
	if err != nil {
		return err
	}

	return nil
}
