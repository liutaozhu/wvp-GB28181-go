package service

import (
	"wvp-pro-go/internal/database"
	"wvp-pro-go/internal/model"
	"wvp-pro-go/internal/utils"

	"go.uber.org/zap"
)

type RegionService struct {
	log *zap.Logger
}

func NewRegionService(log *zap.Logger) *RegionService {
	return &RegionService{log: log}
}

func (s *RegionService) Add(r *model.Region) error {
	return database.DB.Create(r).Error
}

func (s *RegionService) GetPageList(page, count int) (*utils.PageInfo[any], error) {
	var regions []model.Region
	var total int64

	db := database.DB.Model(&model.Region{})
	db.Count(&total)
	if err := db.Offset((page - 1) * count).Limit(count).Order("id ASC").Find(&regions).Error; err != nil {
		return nil, err
	}

	list := make([]any, len(regions))
	for i := range regions {
		list[i] = regions[i]
	}
	return utils.NewPageInfo[any](total, list, page, count), nil
}

func (s *RegionService) GetTreeList() ([]map[string]interface{}, error) {
	var regions []model.Region
	if err := database.DB.Order("id ASC").Find(&regions).Error; err != nil {
		return nil, err
	}

	return buildRegionTree(regions), nil
}

func (s *RegionService) Update(r *model.Region) error {
	return database.DB.Save(r).Error
}

func (s *RegionService) Delete(id uint) error {
	return database.DB.Delete(&model.Region{}, id).Error
}

func (s *RegionService) GetOne(id uint) (*model.Region, error) {
	var r model.Region
	err := database.DB.Where("id = ?", id).First(&r).Error
	return &r, err
}

func (s *RegionService) GetChildren(parentID uint) ([]model.Region, error) {
	var regions []model.Region
	err := database.DB.Where("parent_id = ?", parentID).Find(&regions).Error
	return regions, err
}

func (s *RegionService) GetPath(id uint) ([]model.Region, error) {
	var path []model.Region
	r, err := s.GetOne(id)
	if err != nil {
		return path, err
	}
	path = append(path, *r)
	for r.ParentID > 0 {
		r, err = s.GetOne(r.ParentID)
		if err != nil {
			break
		}
		path = append([]model.Region{*r}, path...)
	}
	return path, nil
}

func buildRegionTree(regions []model.Region) []map[string]interface{} {
	type treeNode struct {
		ID       uint                   `json:"id"`
		Name     string                 `json:"name"`
		PID      uint                   `json:"pid"`
		Children []map[string]interface{} `json:"children,omitempty"`
	}

	nodeMap := make(map[uint]*treeNode)
	var roots []*treeNode

	for _, r := range regions {
		node := &treeNode{ID: r.ID, Name: r.Name, PID: r.ParentID}
		nodeMap[r.ID] = node
		if r.ParentID == 0 {
			roots = append(roots, node)
		}
	}

	for _, r := range regions {
		if r.ParentID > 0 {
			if parent, ok := nodeMap[r.ParentID]; ok {
				parent.Children = append(parent.Children, map[string]interface{}{
					"id":   r.ID,
					"name": r.Name,
					"pid":  r.ParentID,
				})
			}
		}
	}

	result := make([]map[string]interface{}, len(roots))
	for i, r := range roots {
		result[i] = map[string]interface{}{
			"id":   r.ID,
			"name": r.Name,
			"pid":  r.PID,
		}
		if len(r.Children) > 0 {
			result[i]["children"] = r.Children
		}
	}
	return result
}
