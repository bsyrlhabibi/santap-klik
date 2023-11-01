package models

import "gorm.io/gorm"

type SearchResult struct {
	JajananResults []Jajanan
	MakananResults []Makanan
}

type SearchModel struct {
	db *gorm.DB
}

func NewSearchModel(db *gorm.DB) *SearchModel {
	return &SearchModel{db}
}

func (sm *SearchModel) SearchByKeyword(keyword string) (*SearchResult, error) {
	var result SearchResult

	// Lakukan pencarian Jajanan
	if err := sm.db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&result.JajananResults).Error; err != nil {
		return nil, err
	}

	// Lakukan pencarian Makanan
	if err := sm.db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&result.MakananResults).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (sm *SearchModel) SearchJajananByKeyword(keyword string) ([]Jajanan, error) {
	var jajananResults []Jajanan

	if err := sm.db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&jajananResults).Error; err != nil {
		return nil, err
	}

	return jajananResults, nil
}

func (sm *SearchModel) SearchMakananByKeyword(keyword string) ([]Makanan, error) {
	var makananResults []Makanan

	if err := sm.db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&makananResults).Error; err != nil {
		return nil, err
	}

	return makananResults, nil
}
