package repo

const (
	createSurvivor = `INSERT INTO survivor (
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
RETURNING *
`
	getAll = "SELECT * FROM survivor where settlement_id = $1"
)
