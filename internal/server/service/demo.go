package service

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/topvennie/fragtape/internal/database/model"
	"github.com/topvennie/fragtape/internal/database/repository"
	"github.com/topvennie/fragtape/internal/server/dto"
	"github.com/topvennie/fragtape/pkg/storage"
	"github.com/topvennie/fragtape/pkg/utils"
	"go.uber.org/zap"
)

type Demo struct {
	service Service

	demo      repository.Demo
	highlight repository.Highlight
	stat      repository.Stat
	statsDemo repository.StatsDemo
	user      repository.User
}

func (s *Service) NewDemo() *Demo {
	return &Demo{
		service:   *s,
		demo:      *s.repo.NewDemo(),
		highlight: *s.repo.NewHighlight(),
		stat:      *s.repo.NewStat(),
		statsDemo: *s.repo.NewStatsDemo(),
		user:      *s.repo.NewUser(),
	}
}

func (d *Demo) GetAll(ctx context.Context, userID int) ([]dto.Demo, error) {
	demosModel, err := d.demo.GetByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	demoIDs := utils.SliceMap(demosModel, func(d *model.Demo) int { return d.ID })

	statsDemosModel, err := d.statsDemo.GetByDemos(ctx, demoIDs)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	statsDemosMap := make(map[int]*model.StatsDemo)
	for _, s := range statsDemosModel {
		statsDemosMap[s.DemoID] = s
	}

	statsModel, err := d.stat.GetByDemos(ctx, demoIDs)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	statMap := make(map[int][]*model.Stat)
	for _, s := range statsModel {
		demoStats, ok := statMap[s.DemoID]
		if !ok {
			demoStats = []*model.Stat{}
		}

		demoStats = append(demoStats, s)
		statMap[s.DemoID] = demoStats
	}

	highlightsModel, err := d.highlight.GetByDemos(ctx, demoIDs)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	highlightMap := make(map[int]map[int][]*model.Highlight)
	for _, h := range highlightsModel {
		demoHighlights, ok := highlightMap[h.DemoID]
		if !ok {
			highlightMap[h.DemoID] = map[int][]*model.Highlight{}
			demoHighlights = map[int][]*model.Highlight{}
		}
		userHighlights, ok := demoHighlights[h.UserID]
		if !ok {
			userHighlights = []*model.Highlight{}
		}

		userHighlights = append(userHighlights, h)
		highlightMap[h.DemoID][h.UserID] = userHighlights
	}

	users, err := d.user.GetByIDs(ctx, utils.SliceUnique(utils.SliceMap(statsModel, func(s *model.Stat) int { return s.UserID })))
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	userMap := make(map[int]*model.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	demos := make([]dto.Demo, 0, len(demosModel))
	for _, demoModel := range demosModel {
		demo := dto.DemoDTO(demoModel)
		if stat, ok := statsDemosMap[demo.ID]; ok {
			demo.Stats = dto.StatsDemoDTO(stat)
		}

		stats, ok := statMap[demo.ID]
		if !ok {
			continue
		}

		demoHighlights, ok := highlightMap[demo.ID]
		if !ok {
			demoHighlights = map[int][]*model.Highlight{}
		}

		for _, s := range stats {
			user, ok := userMap[s.UserID]
			if !ok {
				continue
			}

			demoPlayer := dto.DemoPlayer{
				User:       dto.UserDTO(user),
				Stat:       dto.StatDTO(s),
				Highlights: utils.SliceMap(demoHighlights[s.UserID], dto.HighlightDTO),
			}

			demo.Players = append(demo.Players, demoPlayer)
		}

		demos = append(demos, demo)
	}

	return demos, nil
}

func (d *Demo) Upload(ctx context.Context, userID int, file []byte) error {
	demo := &model.Demo{
		Source:   model.DemoSourceManual,
		SourceID: strconv.Itoa(userID),
		FileID:   uuid.NewString(),
	}

	return d.service.withRollback(ctx, func(c context.Context) error {
		if err := d.demo.Create(ctx, demo); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if err := d.stat.Create(ctx, &model.Stat{
			DemoID: demo.ID,
			UserID: userID,
		}); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if err := storage.S.Set(demo.FileID, file, 0); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		return nil
	})
}
