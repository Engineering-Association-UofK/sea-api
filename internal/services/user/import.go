package user

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"sea-api/internal/config"
	"sea-api/internal/models"
	"sea-api/internal/utils"
	"sea-api/internal/utils/sheets"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) ImportUsers(eventID int64, file io.Reader) error {
	users, err := sheets.ParseExcelToStructs[models.EventUsersImport](file)
	if err != nil {
		return err
	}

	ids := utils.ExtractField(users, func(u models.EventUsersImport) int64 {
		index, _ := strconv.ParseInt(u.Index, 10, 64)
		return index
	})
	existing, err := s.repo.GetAllByIndices(ids)
	if err != nil {
		return err
	}

	existingMap := utils.FromSlice(existing, func(u models.UserModel) int64 { return u.ID })

	tx, err := s.repo.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, u := range users {
		index, err := strconv.ParseInt(u.Index, 10, 64)
		if err != nil {
			return err
		}
		if _, ok := existingMap[index]; !ok {
			username := sha512.Sum512([]byte(u.NameEn + "|" + config.App.SecretSalt))
			p, _ := generatePasscode(8)
			pass, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
			err = s.repo.Create(&models.UserModel{
				ID:         index,
				UniID:      0,
				Username:   hex.EncodeToString(username[:]),
				NameEn:     u.NameEn,
				NameAr:     u.NameAr,
				Email:      u.Email,
				Phone:      "",
				Department: "",
				Verified:   false,
				Password:   string(pass),
				Status:     models.STATUS_INACTIVE,
				Gender:     models.MALE,
			}, tx)
			if err != nil {
				return err
			}

			err = s.repo.DeleteTempUser(index, tx)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUsersImport(file io.Reader) error {
	mods, err := sheets.ParseExcelToStructs[models.ImportUserUpdate](file)
	if err != nil {
		return err
	}

	users, err := s.repo.GetAll(100, 1)
	if err != nil {
		return err
	}

	modsMap := utils.FromSlice(mods, func(u models.ImportUserUpdate) string { return u.Email })

	for _, u := range users {
		if user, ok := modsMap[u.Email]; ok {
			index, err := strconv.ParseInt(user.Index, 10, 64)
			if err != nil {
				slog.Error("user "+user.Index+" failed to update", "error", err)
			}
			u.ID = index
			u.Phone = user.Phone
			u.Status = models.STATUS_INACTIVE
			s.repo.UpdateWithID(&u, nil)
		}
	}

	return nil
}
