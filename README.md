# SpaceDodge

<p align="center">
  <video src="assets/spacedodge.mov" autoplay loop controls width="700"></video>
</p>

> Survive. Dodge. Repeat.  
> A fast-paced 2D arcade game built with Go + Ebiten.

---

## What is SpaceDodge?

SpaceDodge is a simple but addictive arcade game where you control a spaceship trying to survive an endless asteroid storm.

The rules?

- Asteroids fall from the sky  
- You dodge them  
- You have limited health  
- The game gets harder over time  
- Collisions hurt (a lot)

How long can you survive?

---

## 🎮 Gameplay

- Move Left → `A`
- Move Right → `D`
- Avoid incoming asteroids
- Score increases the longer you survive
- Difficulty scales as your level increases

Pro tip: Panic makes you crash faster.

---

## Built With

- Go
- [Ebiten](https://ebitengine.org/) (2D game engine for Go)
- Sound effects & background music
- Custom starfield background
- Particle explosion effects

---

## How To Run

### 1. Clone the repository

```bash
git clone https://github.com/VincentSamuelPaul/SpaceDodge.git
cd SpaceDodge
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run the Game

```bash
make run
```

## Features

- Animated starfield background
- Rotating asteriods
- Health system
- Sound effects (hit + explosion)
- Particle explosion effect on crash

## Future Improvements

- Power-ups (shield, slow time)
- Shooting asteriods or destruction

## Built by

Vincent Samuel Paul

## If you liked it

- Give it a star
- Beat your high score