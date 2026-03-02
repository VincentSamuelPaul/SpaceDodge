package main

type Player struct {
	x, y float64
	w, h float64
}

type Game struct {
	player                   Player
	stars                    []Star
	obstacleCount            int
	obstacles                []Obstacle
	tiltAngle                float64
	score, health            int
	level                    int
	lastLevelScore           float64
	hitTimer, explosionTimer int
	particles                [][4]float64
}

type Star struct {
	x, y  float64
	speed float64
}

type Obstacle struct {
	x, y, w, h                   float64
	speed, rotation, spin, scale float64
	imgIndex                     int
}
