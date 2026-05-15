package service

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"

	"go.uber.org/zap"
)

type GroupService struct {
	log *zap.Logger
}

func NewGroupService(log *zap.Logger) *GroupService {
	return &GroupService{log: log}
}

func (s *GroupService) Add(g *model.Group) error {
	return database.DB.Create(g).Error
}

func (s *GroupService) GetTreeList() ([]map[string]interface{}, error) {
	var groups []model.Group
	if err := database.DB.Order("id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return buildGroupTree(groups), nil
}

func (s *GroupService) Update(g *model.Group) error {
	return database.DB.Save(g).Error
}

func (s *GroupService) Delete(id uint) error {
	return database.DB.Delete(&model.Group{}, id).Error
}

func (s *GroupService) GetPath(id uint) ([]model.Group, error) {
	var path []model.Group
	g, err := s.GetOne(id)
	if err != nil {
		return path, err
	}
	path = append(path, *g)
	for g.ParentID > 0 {
		g, err = s.GetOne(g.ParentID)
		if err != nil {
			break
		}
		path = append([]model.Group{*g}, path...)
	}
	return path, nil
}

func (s *GroupService) GetOne(id uint) (*model.Group, error) {
	var g model.Group
	err := database.DB.Where("id = ?", id).First(&g).Error
	return &g, err
}

func buildGroupTree(groups []model.Group) []map[string]interface{} {
	type treeNode struct {
		ID       uint
		Name     string
		PID      uint
		Children []map[string]interface{}
	}

	nodeMap := make(map[uint]*treeNode)
	var roots []*treeNode

	for _, g := range groups {
		node := &treeNode{ID: g.ID, Name: g.Name, PID: g.ParentID}
		nodeMap[g.ID] = node
		if g.ParentID == 0 {
			roots = append(roots, node)
		}
	}

	for _, g := range groups {
		if g.ParentID > 0 {
			if parent, ok := nodeMap[g.ParentID]; ok {
				parent.Children = append(parent.Children, map[string]interface{}{
					"id":   g.ID,
					"name": g.Name,
					"pid":  g.ParentID,
				})
			}
		}
	}

	result := make([]map[string]interface{}, len(roots))
	for i, g := range roots {
		result[i] = map[string]interface{}{
			"id":   g.ID,
			"name": g.Name,
			"pid":  g.PID,
		}
		if len(g.Children) > 0 {
			result[i]["children"] = g.Children
		}
	}
	return result
}
