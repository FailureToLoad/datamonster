package survivor

import (
	"encoding/json"
	storeMocks "github.com/failuretoload/datamonster/store/mocks"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
)

type SurvivorApiTestSuite struct {
	suite.Suite
	target *Controller
	db     *storeMocks.MockConnection
	router *chi.Mux
}

func (suite *SurvivorApiTestSuite) SetupTest() {
	suite.db = &storeMocks.MockConnection{}
	suite.target = NewController(suite.db)
	suite.router = chi.NewRouter()
	suite.target.RegisterRoutes(suite.router)
}

func (suite *SurvivorApiTestSuite) Test_GetSurvivors_ReturnsSurvivorList() {
	rows := storeMocks.MockRows{
		Rows: []pgx.Row{
			&SurvivorRow{
				Id:               1,
				Settlement:       1,
				Name:             "Zach",
				Born:             1,
				Gender:           "M",
				Status:           "alive",
				HuntXp:           1,
				Survival:         1,
				Movement:         1,
				Accuracy:         1,
				Strength:         1,
				Evasion:          1,
				Luck:             1,
				Speed:            1,
				Insanity:         1,
				SystemicPressure: 1,
				Torment:          1,
				Lumi:             1,
				Understanding:    1,
				Courage:          1,
			},
			&SurvivorRow{
				Id:               2,
				Settlement:       1,
				Name:             "Lucy",
				Born:             1,
				Gender:           "M",
				Status:           "alive",
				HuntXp:           1,
				Survival:         1,
				Movement:         1,
				Accuracy:         1,
				Strength:         1,
				Evasion:          1,
				Luck:             1,
				Speed:            1,
				Insanity:         1,
				SystemicPressure: 1,
				Torment:          1,
				Lumi:             1,
				Understanding:    1,
				Courage:          1,
			},
		},
	}
	suite.db.SetRows(&rows)
	req := httptest.NewRequest("GET", "/settlement/1/survivor", nil)
	ctx := req.Context()
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	suite.Equal(200, resp.StatusCode, "200 response should be returned")
	body, _ := io.ReadAll(resp.Body)
	dtoList := []SurvivorDTO{}
	json.Unmarshal(body, &dtoList)
	count := len(dtoList)
	suite.Equal(2, count, "2 settlements should be returned")
}

func TestSurvivorApiTestSuite(t *testing.T) {
	suite.Run(t, new(SurvivorApiTestSuite))
}

type SurvivorRow struct {
	Id               int
	Settlement       int
	Name             string
	Born             int
	Gender           string
	Status           string
	HuntXp           int
	Survival         int
	Movement         int
	Accuracy         int
	Strength         int
	Evasion          int
	Luck             int
	Speed            int
	Insanity         int
	SystemicPressure int
	Torment          int
	Lumi             int
	Courage          int
	Understanding    int
}

func (s *SurvivorRow) Scan(dest ...interface{}) error {
	id := dest[0].(*int)
	settlement := dest[1].(*int)
	name := dest[2].(*string)
	born := dest[3].(*int)
	gender := dest[4].(*string)
	status := dest[5].(*string)
	huntXp := dest[6].(*int)
	survival := dest[7].(*int)
	movement := dest[8].(*int)
	accuracy := dest[9].(*int)
	strength := dest[10].(*int)
	evasion := dest[11].(*int)
	luck := dest[12].(*int)
	speed := dest[13].(*int)
	insanity := dest[14].(*int)
	systemicPressure := dest[15].(*int)
	torment := dest[16].(*int)
	lumi := dest[17].(*int)
	courage := dest[18].(*int)
	understanding := dest[19].(*int)

	*id = s.Id
	*settlement = s.Settlement
	*name = s.Name
	*born = s.Born
	*gender = s.Gender
	*huntXp = s.HuntXp
	*status = s.Status
	*survival = s.Survival
	*movement = s.Movement
	*accuracy = s.Accuracy
	*strength = s.Strength
	*evasion = s.Evasion
	*luck = s.Luck
	*speed = s.Speed
	*insanity = s.Insanity
	*systemicPressure = s.SystemicPressure
	*torment = s.Torment
	*lumi = s.Lumi
	*courage = s.Courage
	*understanding = s.Understanding
	return nil
}
