package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 640
	screenHeight = 640

	playerSpeed = 10
)

var (
	shipImage    *ebiten.Image
	asteroidImgs []*ebiten.Image

	normalFont font.Face

	audioCtx       *audio.Context
	bgm            *audio.Player
	hitSound       []byte
	explosionSound []byte
)

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyA) && g.player.x > 20 {
		g.player.x -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.player.x+g.player.w < screenWidth-20 {
		g.player.x += playerSpeed
	}

	for i := range g.stars {
		g.stars[i].y += g.stars[i].speed
		if g.stars[i].y > screenHeight {
			g.stars[i].y = 0
			g.stars[i].x = float64(rand.Intn(screenWidth))
		}
	}

	for i := range g.obstacles {
		g.obstacles[i].y += g.obstacles[i].speed
		g.obstacles[i].rotation += g.obstacles[i].spin

		if g.obstacles[i].y > screenHeight {
			g.obstacles[i] = generateObstacle()
		}
	}

	for i := range g.obstacles {
		if collision(g.player, g.obstacles[i]) && g.hitTimer == 0 && g.explosionTimer == 0 {
			g.health--
			playSound(hitSound)
			g.obstacles[i].y = screenHeight + 100
			g.hitTimer = 60
			if g.health <= 0 {
				g.explosionTimer = 60
				playSound(explosionSound)
				g.particles = nil

				for i := 0; i < 50; i++ {
					angle := rand.Float64() * 2 * 3.1415
					speed := 2 + rand.Float64()*2
					dx := speed * math.Cos(angle)
					dy := speed * math.Sin(angle)
					g.particles = append(g.particles, [4]float64{g.player.x + 20, g.player.y + 20, dx, dy})
				}
			}
		}
	}

	if g.score-int(g.lastLevelScore) >= 1800 {
		g.level++
		g.obstacleCount++
		g.lastLevelScore = float64(g.score)
		g.obstacles = append(g.obstacles, generateObstacle())
	}

	if g.explosionTimer > 0 {
		g.explosionTimer--
		for i := range g.particles {
			g.particles[i][0] += g.particles[i][2]
			g.particles[i][1] += g.particles[i][3]
		}

		if g.explosionTimer == 0 {
			g.resetGame()
		}
	}

	g.score++

	// g.player.y = screenHeight - g.player.h - 20

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, s := range g.stars {
		brightness := uint8(155 + rand.Intn(100))
		screen.Set(int(s.x), int(s.y), color.RGBA{brightness, brightness, brightness, 255})
	}

	for _, o := range g.obstacles {
		// ebitenutil.DrawRect(screen, o.x, o.y, o.w, o.h, color.RGBA{100, 100, 100, 255})
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(o.scale, o.scale)
		opts.GeoM.Rotate(o.rotation)
		img := asteroidImgs[o.imgIndex]
		w, h := img.Size()
		opts.GeoM.Translate(-float64(w)*o.scale/2, -float64(h)*o.scale/2)
		opts.GeoM.Translate(o.x+o.w/2, o.y+o.h/2)
		screen.DrawImage(img, opts)
	}

	// ebitenutil.DrawRect(screen, g.player.x, g.player.y, g.player.w, g.player.h, color.RGBA{0, 255, 255, 255})
	opts := &ebiten.DrawImageOptions{}
	scale := 0.7
	opts.GeoM.Scale(scale, scale)
	w, h := shipImage.Size()
	opts.GeoM.Translate(-float64(w)*scale/2, -float64(h)*scale/2)
	shipX := g.player.x + g.player.w/2
	shipY := screenHeight - float64(h)*scale/2 - 60
	opts.GeoM.Translate(shipX, shipY)

	if g.explosionTimer > 0 {
		alpha := uint8(255 * g.explosionTimer / 60)
		for _, p := range g.particles {
			c := color.RGBA{
				R: 255,
				G: uint8(rand.Intn(150)),
				B: 0,
				A: alpha,
			}
			screen.Set(int(p[0]), int(p[1]), c)
		}
		return
	}
	if g.hitTimer > 0 {
		if (g.hitTimer/5)%2 == 0 {
			screen.DrawImage(shipImage, opts)
		}
		g.hitTimer--
	} else {
		screen.DrawImage(shipImage, opts)
	}

	// screen.DrawImage(shipImage, opts)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d\nHealth: %d\nLevel: %d", int(g.score/60), g.health, g.level))

	// msg := fmt.Sprintf("Score: %d\nHealth: %d\nLevel: %d", int(g.score/60), g.health, g.level)
	// text.Draw(screen, msg, normalFont, 20, 40, color.RGBA{255, 255, 0, 255})

	text.Draw(screen, fmt.Sprintf("Score: %d", int(g.score/60)), normalFont, 20, 40, color.RGBA{255, 255, 0, 255})
	text.Draw(screen, fmt.Sprintf("Health: %d", g.health), normalFont, 20, 70, color.RGBA{0, 255, 0, 255})
	text.Draw(screen, fmt.Sprintf("Level: %d", g.level), normalFont, 20, 100, color.RGBA{0, 200, 255, 255})

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func generateStars() []Star {

	stars := make([]Star, 100)

	for i := range stars {
		stars[i] = Star{
			x:     float64(rand.Intn(screenWidth)),
			y:     float64(rand.Intn(screenHeight)),
			speed: 1 + rand.Float64(),
		}
	}
	return stars
}

func (g *Game) resetGame() {
	g.health = 3
	g.score = 0
	g.level = 1
	g.obstacles = generateObstacles(2)
}

func generateObstacle() Obstacle {
	obstacle := Obstacle{
		x:        float64(rand.Intn(screenWidth)),
		y:        0,
		speed:    5 + rand.Float64(),
		w:        float64(20 + rand.Intn(30)),
		h:        float64(10 + rand.Intn(30)),
		spin:     (rand.Float64() - 0.5) * 0.04,
		scale:    0.5 + rand.Float64(),
		imgIndex: rand.Intn(len(asteroidImgs)),
	}

	return obstacle
}

func generateObstacles(n int) []Obstacle {
	obstacles := make([]Obstacle, n)
	for i := range obstacles {
		obstacles[i] = generateObstacle()
		obstacles[i].y = float64(-rand.Intn(screenHeight))
	}
	return obstacles
}

func playSound(data []byte) {
	if data == nil {
		return
	}

	s, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(data))
	if err != nil {
		return
	}
	p, err := audio.NewPlayer(audioCtx, s)
	if err != nil {
		return
	}
	p.Play()
}

func collision(p Player, o Obstacle) bool {
	return p.x < o.x+o.w && p.x+p.w > o.x && p.y < o.y+o.h && p.y+p.h > o.y
}

func main() {

	asteroidFiles := []string{
		"./assets/asteroid1.png",
		"./assets/asteroid2.png",
	}

	for _, file := range asteroidFiles {
		img, _, err := ebitenutil.NewImageFromFile(file)
		if err != nil {
			log.Fatal(err)
		}
		asteroidImgs = append(asteroidImgs, img)
	}

	ttfData, err := os.ReadFile("./assets/fonts.ttf")
	if err != nil {
		log.Fatal(err)
	}

	ttf, err := opentype.Parse(ttfData)
	if err != nil {
		log.Fatal(err)
	}

	normalFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err != nil {
		log.Fatal(err)
	}

	audioCtx = audio.NewContext(44100)

	bgmFile, err := os.Open("./assets/music.mp3")
	if err != nil {
		log.Fatal("music load:", err)
	}
	defer bgmFile.Close()

	bgmStream, err := mp3.DecodeWithSampleRate(44100, bgmFile)
	if err != nil {
		log.Fatal("decode music:", err)
	}

	bgm, err = audio.NewPlayer(audioCtx, audio.NewInfiniteLoop(bgmStream, bgmStream.Length()))
	if err != nil {
		log.Fatal(err)
	}
	bgm.Play()

	hitFile, _ := os.ReadFile("./assets/hit.wav")
	explFile, _ := os.ReadFile("./assets/explosion.wav")

	hitSound = hitFile
	explosionSound = explFile

	game := Game{
		player: Player{
			x: screenWidth/2 - 20,
			y: screenHeight - 100,
			w: 40,
			h: 40,
		},
		stars:          generateStars(),
		obstacles:      generateObstacles(2),
		obstacleCount:  2,
		score:          0,
		health:         3,
		level:          1,
		lastLevelScore: 0,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Space Dodge")

	shipImage, _, err = ebitenutil.NewImageFromFile("./assets/ship.png")
	if err != nil {
		log.Fatal(err)
	}
	w, h := shipImage.Size()
	game.player.w = float64(w)
	game.player.h = float64(h)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
