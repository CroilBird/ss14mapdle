package models

import (
	"errors"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"image"
	"io"
	"log/slog"
	"os"
	"path"
	"ss14mapdle/config"
	"ss14mapdle/util"
)

type BoundingBox struct {
	ID uint64 `json:"id"`
	X  int64  `json:"x"`
	Y  int64  `json:"y"`
	W  int64  `json:"w"`
	H  int64  `json:"h"`
}

type Map struct {
	ID                uint64        `json:"-"`
	Name              string        `json:"name"`
	Path              string        `json:"path"`
	Index             int           `json:"index"`
	Width             uint64        `json:"-"`
	Height            uint64        `json:"-"`
	WithinBoundsBoxes []BoundingBox `json:"within_bounds_boxes" gorm:"many2many:map_bounding_boxes"`
}

func init() {
	Models = append(Models, Map{})
}

func GetMap(db *gorm.DB, mapName string) (*Map, error) {
	var selectedMap Map

	err := db.Table("maps").
		Where("name = ?", mapName).
		Find(&selectedMap).
		Error

	if err != nil {
		return nil, err
	}

	return &selectedMap, nil
}

func GetRandomMap(db *gorm.DB) (*Map, error) {
	var randomMap Map

	err := db.Table("maps").
		Order("RANDOM()").
		Last(&randomMap).
		Error

	if err != nil {
		return nil, err
	}

	return &randomMap, nil
}

func LoadMapJson(db *gorm.DB) error {
	var err error
	var file *os.File

	mapBasePath, err := config.GetConfig(config.EnvMapBasePath)

	if err != nil {
		return err
	}

	mapInfoPath := path.Join(mapBasePath, "maps.json")

	defer func() {
		deferErr := file.Close()
		if deferErr != nil {
			panic("this should not happen bye")
		}
	}()

	if file, err = os.Open(mapInfoPath); err != nil {
		slog.Error("[util] Could not open entry file", "path", mapInfoPath)
		return err
	}

	var bytes []byte
	if bytes, err = io.ReadAll(file); err != nil {
		slog.Error("[util] Could not read entry file", "path", mapInfoPath)
		return err
	}

	var maps []Map

	err = json.Unmarshal(bytes, &maps)

	if err != nil {
		slog.Error("[util] Could not unmarshal map entry data", "err", err)
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, mapInfo := range maps {
			var count int64
			transactionErr := tx.
				Table("maps").
				Where("name = ? AND `index` = ?", mapInfo.Name, mapInfo.Index).
				Count(&count).
				Error

			if transactionErr != nil {
				return transactionErr
			}

			if count != 0 {
				continue
			}

			imgPtr, err := util.ReadMap(mapInfo.Path, int(mapInfo.Index))

			if err != nil {
				return err
			}

			img := *imgPtr

			type imageWithBounds interface {
				Bounds() image.Rectangle
			}

			subImg, ok := img.(imageWithBounds)
			if !ok {
				return errors.New("could not do some bullshit idk")
			}

			imgSize := subImg.Bounds().Size()

			mapInfo.Width = uint64(imgSize.X)
			mapInfo.Height = uint64(imgSize.Y)

			transactionErr = tx.
				Create(&mapInfo).
				Error

			if transactionErr != nil {
				return transactionErr
			}
		}

		return nil
	})

	return err
}
