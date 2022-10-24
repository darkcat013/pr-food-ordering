package config

import "time"

const LOGS_ENABLED = true
const PORT = "8088"

const TIME_UNIT = time.Millisecond * TIME_UNIT_COEFF
const TIME_UNIT_COEFF = 100
