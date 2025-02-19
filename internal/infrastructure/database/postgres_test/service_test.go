package postgres_test

import (
	"context"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"
)

func (suite *PostgresTestSuite) TestServiceQueries() {
	ctx := context.Background()
	// (3, 'Аренда книги', 50.00, 'Почасовая аренда редких книг из библиотеки.');

	exp := &orgdto.AddServiceReq{
		OrgID: 3,
		ServiceInfo: entity.Service{
			Name:        "Аренда книги",
			Cost:        50.0,
			Description: "Почасовая аренда редких книг из библиотеки",
		},
	}
	serviceID, err := suite.db.ServiceAdd(ctx, orgmap.AddServiceToModel(exp))
	suite.NoError(err)
	suite.NotZero(serviceID)

	service, err := suite.db.Service(ctx, serviceID, exp.OrgID)
	suite.NoError(err)
	suite.NotNil(service)
	suite.Equal(&exp.ServiceInfo, orgmap.ServiceToEntity(service))

	expected := &orgdto.UpdateServiceReq{
		ServiceID: serviceID,
		OrgID:     3,
		ServiceInfo: entity.Service{
			Name:        "Аренда книги",
			Cost:        9999.0,
			Description: "ТЕСТИРОВАНИЕ",
		},
	}

	suite.NoError(suite.db.ServiceUpdate(ctx, orgmap.UpdateService(expected)))

	serviceList, found, err := suite.db.ServiceList(ctx, expected.OrgID, 5, 0)
	suite.NoError(err)
	suite.NotNil(serviceList)
	suite.NotZero(found)

	for i := range serviceList {
		service := orgmap.ServiceToDTO(serviceList[i])
		if expected.ServiceID == service.ServiceID {
			suite.Equal(&expected.ServiceInfo, service.ServiceInfo)
		}
	}

	suite.NoError(suite.db.ServiceDelete(ctx, expected.ServiceID, expected.OrgID))
	serviceList, found, err = suite.db.ServiceList(ctx, expected.OrgID, 5, 0)
	suite.NoError(err)
	suite.NotNil(serviceList)
	suite.Zero(found)

	for i := range serviceList {
		suite.NotEqual(&expected.ServiceInfo, orgmap.ServiceToDTO(serviceList[i]).ServiceInfo)
	}

}
