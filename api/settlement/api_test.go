package settlement

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SettlementApiTestSuite struct {
	suite.Suite
	target *Controller
}

func (suite *SettlementApiTestSuite) SetupTest() {
	suite.target = NewController(NewRepo(nil))
}

func (suite *SettlementApiTestSuite) TestSomething() {
	suite.target.getSettlements(nil, nil)
	assert.Equal(suite.T(), 1, 1, "they should be equal")
}

func TestSettlementApiTestSuite(t *testing.T) {
	suite.Run(t, new(SettlementApiTestSuite))
}
