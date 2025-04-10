package postgres

import (
	"context"
	"time"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"

	"github.com/google/uuid"
)

func (suite *PostgresTestSuite) TestOrganizationQueries() {
	ctx := context.Background()
	expOrg := entity.OrgInfo{
		Name:        "TestCompany",
		Address:     "Testovaya street",
		Type:        "TestType",
		Telephone:   "79998887744",
		City:        "Moscow",
		Coordinates: entity.Coordinates{Lat: 56.1252, Long: 38.6173},
	}
	registerInfo := &authdto.OrgRegisterReq{
		UUID:        uuid.New().String(),
		Credentials: authdto.Credentials{Email: "email@test.ru", Password: "testpasswd"},
		OrgInfo:     expOrg,
	}
	orgID, err := suite.db.OrgSave(ctx, orgmap.RegisterReqToModel(registerInfo))
	suite.NoError(err)
	suite.Greater(orgID, 0)

	actOrg, err := suite.db.OrgByID(ctx, orgID)
	suite.NoError(err)
	suite.NotNil(actOrg)
	suite.Equal(orgID, orgmap.OrganizationToDTO(actOrg, time.Now().Location()).OrgID)
	suite.Equal(&expOrg, orgmap.OrganizationToDTO(actOrg, time.Now().Location()).Info)

	expOrg.Address = "another test street"
	expOrg.Type = "testtype"

	updateInfo := &orgdto.OrgUpdateReq{
		OrgID:   orgID,
		OrgInfo: expOrg,
	}

	suite.Require().NoError(suite.db.OrgUpdate(ctx, orgmap.OrgUpdateToModel(updateInfo)))

	params := &general.SearchReq{
		Page:       1,
		Limit:      5,
		Name:       expOrg.Name,
		IsRateSort: true,
	}
	foundOrgs, err := suite.db.OrgsBySearch(ctx, orgmap.SearchToModel(params))
	suite.NoError(err, "without sort")
	suite.Greater(foundOrgs.Found, 0)
	suite.NotNil(foundOrgs)
	suite.NotNil(foundOrgs.Data)
	for i := range foundOrgs.Data {
		org := orgmap.OrgsBySearchToDTO(foundOrgs.Data[i], time.Now().Location())
		if org.OrgID == orgID {
			suite.Equal(expOrg.Address, org.Address)
			suite.Equal(expOrg.Type, org.Type)
		}
	}
	params.IsRateSort = true
	foundOrgs, err = suite.db.OrgsBySearch(ctx, orgmap.SearchToModel(params))
	suite.NoError(err, "rate sort")
	suite.Greater(foundOrgs.Found, 0)
	suite.NotNil(foundOrgs)
	suite.NotNil(foundOrgs.Data)

	params.IsRateSort = false
	params.IsNameSort = true
	foundOrgs, err = suite.db.OrgsBySearch(ctx, orgmap.SearchToModel(params))
	suite.NoError(err, "name sort")
	suite.Greater(foundOrgs.Found, 0)
	suite.NotNil(foundOrgs)
	suite.NotNil(foundOrgs.Data)

	params.IsRateSort = true
	params.IsNameSort = true
	foundOrgs, err = suite.db.OrgsBySearch(ctx, orgmap.SearchToModel(params))
	suite.NoError(err, "rate & name sort")
	suite.Greater(foundOrgs.Found, 0)
	suite.NotNil(foundOrgs)
	suite.NotNil(foundOrgs.Data)

	expOrg.Name = "WELL KNOWN COMPANY"
	expOrg.Type = "Mega Corporation"

	updateInfo = &orgdto.OrgUpdateReq{
		OrgID:   orgID,
		OrgInfo: expOrg,
	}

	suite.Require().NoError(suite.db.OrgUpdate(ctx, orgmap.OrgUpdateToModel(updateInfo)))
	areaParams := &general.OrgAreaReq{
		LeftLowerCorner:  entity.Coordinates{Lat: 50.0, Long: 30.0},
		RightUpperCorner: entity.Coordinates{Lat: 60.0, Long: 40.0},
	}
	areaOrgs, err := suite.db.OrgsInArea(ctx, orgmap.AreaToModel(areaParams))
	suite.NoError(err)
	suite.NotNil(areaOrgs)

	resp := &general.OrgAreaResp{
		Found: len(areaOrgs),
		Orgs:  make([]*entity.MapOrgInfo, 0, len(areaOrgs)),
	}
	suite.Greater(resp.Found, 0)
	for _, v := range areaOrgs {
		resp.Orgs = append(resp.Orgs, orgmap.OrgSummaryToDTO(v))
	}
	for i := range resp.Orgs {
		if resp.Orgs[i].OrgID == orgID {
			suite.Equal(expOrg.Name, resp.Orgs[i].Name)
			suite.Equal(expOrg.Type, resp.Orgs[i].Type)
		}
	}
	suite.Require().NoError(suite.db.OrgSoftDelete(ctx, orgID))
	suite.NoError(suite.db.OrgDelete(ctx, orgID))
}
