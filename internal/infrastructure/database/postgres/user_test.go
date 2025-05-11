//go:build integration

package postgres

import (
	"context"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/authdto"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/infrastructure/mapper/usermap"

	"github.com/google/uuid"
)

func (suite *PostgresTestSuite) TestUserQueries() {
	ctx := context.Background()

	expUser := entity.User{
		UUID:      "",
		FirstName: "test",
		LastName:  "testovich",
		Telephone: "+79876543210",
		City:      "Moscow",
		Email:     "test@test.ru",
	}
	registerInfo := &authdto.UserRegisterReq{
		UUID:        expUser.UUID,
		Credentials: authdto.Credentials{Email: "test@test.ru", Password: "testpasswd"},
		User:        expUser,
	}
	id, err := suite.db.UserSave(ctx, usermap.UserRegisterToModel(registerInfo))
	suite.NoError(err)
	suite.Greater(id, 0)

	expUser.UserID = id

	userDB, err := suite.db.UserByID(ctx, expUser.UserID)
	suite.NoError(err)
	suite.Require().NotNil(userDB)
	actUser := usermap.UserInfoToDTO(userDB)
	suite.Equal(&expUser, actUser)

	expUser.FirstName = "megatest"
	newUUID, err := uuid.NewRandom()
	suite.NoError(err)
	expUser.UUID = newUUID.String()
	updateUser := &userdto.UserUpdateReq{
		UserID:    expUser.UserID,
		FirstName: expUser.FirstName,
		LastName:  expUser.LastName,
		Telephone: expUser.Telephone,
		City:      expUser.City,
	}
	suite.NoError(suite.db.UserUpdate(ctx, usermap.UserUpdateToModel(updateUser)))
	suite.NoError(suite.db.UserSetUUID(ctx, expUser.UserID, expUser.UUID))

	userDB, err = suite.db.UserByID(ctx, expUser.UserID)
	suite.NoError(err)
	suite.Require().NotNil(userDB)
	actUser = usermap.UserInfoToDTO(userDB)
	suite.Equal(&expUser, actUser)

	suite.NoError(suite.db.UserDeleteURL(ctx, expUser.UserID, expUser.UUID))
	userUUID, err := suite.db.UserUUID(ctx, expUser.UserID)
	suite.NoError(err)
	suite.Empty(userUUID)
	suite.NotEqual(expUser.UUID, userUUID)

	suite.Require().NoError(suite.db.UserSoftDelete(ctx, expUser.UserID))
	suite.NoError(suite.db.UserDelete(ctx, expUser.UserID))
}
