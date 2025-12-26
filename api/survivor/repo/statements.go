package repo

import "fmt"

const (
	insertSurvivorBase = `INSERT INTO survivor (
	settlement_id,
	name,
	birth,
	gender,
	hunt_xp,
	survival,
	movement,
	accuracy,
	strength,
	evasion,
	luck,
	speed,
	insanity,
	systemic_pressure,
	torment,
	lumi,
	courage,
	understanding
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11,
	$12,
	$13,
	$14,
	$15,
	$16,
	$17,
	$18
)
`
	createSurvivor = insertSurvivorBase + `RETURNING *`
	getAll         = "SELECT * FROM survivor where settlement_id = $1"
)

var upsertSurvivor = fmt.Sprintf(`%s
ON CONFLICT (settlement_id, name)
DO UPDATE SET
	birth = EXCLUDED.birth,
	gender = EXCLUDED.gender,
	hunt_xp = EXCLUDED.hunt_xp,
	survival = EXCLUDED.survival,
	movement = EXCLUDED.movement,
	accuracy = EXCLUDED.accuracy,
	strength = EXCLUDED.strength,
	evasion = EXCLUDED.evasion,
	luck = EXCLUDED.luck,
	speed = EXCLUDED.speed,
	insanity = EXCLUDED.insanity,
	systemic_pressure = EXCLUDED.systemic_pressure,
	torment = EXCLUDED.torment,
	lumi = EXCLUDED.lumi,
	courage = EXCLUDED.courage,
	understanding = EXCLUDED.understanding
RETURNING *
`, insertSurvivorBase)
