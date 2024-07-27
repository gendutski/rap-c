package unitusecase

import (
	"context"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/repository/contract"
	usecasecontract "rap-c/app/usecase/contract"
	"rap-c/config"
)

func NewUsecase(cfg *config.Config, unitRepo contract.UnitRepository) usecasecontract.UnitUsecase {
	return &usecase{cfg, unitRepo}
}

type usecase struct {
	cfg      *config.Config
	unitRepo contract.UnitRepository
}

func (uc *usecase) Create(ctx context.Context, payload *payloadentity.CreateUnitPayload, author *databaseentity.User) (*databaseentity.Unit, error) {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return nil, err
	}

	// create unit
	unit := databaseentity.Unit{
		Name:      payload.Name,
		CreatedBy: author.ID,
	}
	err = uc.unitRepo.Create(ctx, &unit)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (uc *usecase) GetUnitList(ctx context.Context, req *payloadentity.GetUnitListRequest) ([]*databaseentity.Unit, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	return uc.unitRepo.GetUnitsByRequest(ctx, req)
}

func (uc *usecase) GetTotalUnitList(ctx context.Context, req *payloadentity.GetUnitListRequest) (int64, error) {
	return uc.unitRepo.GetTotalUnitsByRequest(ctx, req)
}

func (uc *usecase) Delete(ctx context.Context, payload *payloadentity.DeleteUnitPayload) error {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return err
	}

	// get unit
	unit, err := uc.unitRepo.GetUnitByName(ctx, payload.Name)
	if err != nil {
		return err
	}

	// delete unit
	return uc.unitRepo.Delete(ctx, unit)
}
