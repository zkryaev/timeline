package postgres

import (
	"context"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure/mapper/orgmap"
)

func (suite *PostgresTestSuite) TestWorkerQueries() {
	ctx := context.Background()

	expWorker := entity.Worker{
		FirstName:       "Donald",
		LastName:        "Trump",
		Position:        "President",
		Degree:          "PhD Politic",
		SessionDuration: 10,
	}
	org, err := suite.db.OrgByID(ctx, 1)
	suite.Require().NoError(err)
	suite.Require().NotNil(org)

	addreq := &orgdto.AddWorkerReq{
		OrgID:      org.OrgID,
		WorkerInfo: expWorker,
	}
	workerID, err := suite.db.WorkerAdd(ctx, orgmap.AddWorkerToModel(addreq))
	suite.Require().NoError(err)
	suite.Require().Greater(workerID, 0)

	dbWorker, err := suite.db.Worker(ctx, workerID, org.OrgID)
	suite.NoError(err)
	suite.Require().NotNil(dbWorker)
	actWorker := orgmap.WorkerToDTO(dbWorker)
	suite.Equal(&expWorker, actWorker.WorkerInfo)

	expWorker.FirstName = "Putin"
	expWorker.LastName = "Vladimirovich"
	expWorker.Position = "Imperor of the Mankind!"
	expWorker.Degree = "Doesn't need"
	updreq := &orgdto.UpdateWorkerReq{
		WorkerID:   workerID,
		OrgID:      org.OrgID,
		WorkerInfo: expWorker,
	}
	suite.Require().NoError(suite.db.WorkerUpdate(ctx, orgmap.UpdateWorkerToModel(updreq)))

	dbWorkers, found, err := suite.db.WorkerList(ctx, org.OrgID, 10, 0)
	suite.NoError(err)
	suite.Greater(found, 0)
	suite.Require().NotNil(dbWorkers)
	for i := range dbWorkers {
		actworker := orgmap.WorkerToDTO(dbWorkers[i])
		if actworker.WorkerID == workerID {
			suite.Equal(&expWorker, actworker.WorkerInfo)
		}
	}

	services, found, err := suite.db.ServiceList(ctx, org.OrgID, 10, 0)
	suite.NoError(err)
	suite.Greater(found, 0)
	suite.Require().NotNil(services)
	service := orgmap.ServiceToDTO(services[0])

	assreq := &orgdto.AssignWorkerReq{
		ServiceID: service.ServiceID,
		OrgID:     org.OrgID,
		WorkerID:  workerID,
	}
	suite.Require().NoError(suite.db.WorkerAssignService(ctx, orgmap.AssignWorkerToModel(assreq)))

	serviceWorkers, err := suite.db.ServiceWorkerList(ctx, service.ServiceID, org.OrgID)
	suite.NoError(err)
	suite.Require().NotNil(serviceWorkers)
	for i := range serviceWorkers {
		worker := orgmap.WorkerToDTO(dbWorkers[i])
		if worker.WorkerID == workerID {
			suite.Equal(&expWorker, worker.WorkerInfo)
		}
	}

	suite.Require().NoError(suite.db.WorkerUnAssignService(ctx, orgmap.AssignWorkerToModel(assreq)))
	serviceWorkers, err = suite.db.ServiceWorkerList(ctx, service.ServiceID, org.OrgID)
	suite.NoError(err)
	suite.Require().NotNil(serviceWorkers)
	for i := range serviceWorkers {
		suite.NotEqual(workerID, serviceWorkers[i].WorkerID)
	}
	suite.Require().NoError(suite.db.WorkerSoftDelete(ctx, workerID, org.OrgID))
	suite.NoError(suite.db.WorkerDelete(ctx, workerID, org.OrgID))
}
