//go:build integration
package postgres

import (
	"context"
	"time"
	"timeline/internal/infrastructure/models"
)

func (suite *PostgresTestSuite) TestSaveVerifyCodeUser() {
	ctx := context.Background()
	user, err := suite.db.UserByID(ctx, 1)
	suite.Require().NoError(err)
	userCode := &models.CodeInfo{
		ID:    user.UserID,
		Code:  "0000",
		IsOrg: false,
	}
	err = suite.db.SaveVerifyCode(ctx, userCode)
	suite.NoError(err, user.FirstName, user.UserID)

	exp, err := suite.db.VerifyCode(ctx, userCode)
	suite.NoError(err)
	suite.NotEqual(time.Time{}, exp)

	err = suite.db.ActivateAccount(ctx, userCode.ID, userCode.IsOrg)
	suite.NoError(err)

	expAcc, err := suite.db.AccountExpiration(ctx, user.Email, userCode.IsOrg)
	suite.NoError(err)
	suite.NotEqual(models.ExpInfo{}, expAcc)
}

func (suite *PostgresTestSuite) TestSaveVerifyCodeOrg() {
	ctx := context.Background()
	org, err := suite.db.OrgByID(ctx, 1)
	suite.Require().NoError(err)
	orgCode := &models.CodeInfo{
		ID:    org.OrgID,
		Code:  "0000",
		IsOrg: true,
	}
	err = suite.db.SaveVerifyCode(ctx, orgCode)
	suite.NoError(err)

	exp, err := suite.db.VerifyCode(ctx, orgCode)
	suite.NoError(err)
	suite.NotEqual(time.Time{}, exp)

	err = suite.db.ActivateAccount(ctx, orgCode.ID, orgCode.IsOrg)
	suite.NoError(err)

	expAcc, err := suite.db.AccountExpiration(ctx, org.Email, orgCode.IsOrg)
	suite.NoError(err)
	suite.NotEqual(models.ExpInfo{}, expAcc)
}
