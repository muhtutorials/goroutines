package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	id       int
	position Vector2D
	velocity Vector2D
}

func (b *Boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}
	return 0
}

func (b *Boid) calcAcceleration() Vector2D {
	// boid's view square
	upperRight, lowerLeft := b.position.AddVal(viewRadius), b.position.AddVal(-viewRadius)
	avgPosition, avgVelocity, separation := Vector2D{x: 0, y: 0}, Vector2D{x: 0, y: 0}, Vector2D{x: 0, y: 0}
	count := 0.0

	rWlock.RLock()
	// min and max remove pixels outside screen
	for i := math.Max(lowerLeft.x, 0); i <= math.Min(upperRight.x, screenWidth); i++ {
		for j := math.Max(lowerLeft.y, 0); j <= math.Min(upperRight.y, screenHeight); j++ {
			if otherBoidId := boidMap[int(i)][int(j)]; otherBoidId != -1 && otherBoidId != b.id {
				// square is bigger than circle, so we exclude boids that are not in the circle
				if distance := boids[otherBoidId].position.Distance(b.position); distance < viewRadius {
					count++
					avgPosition = avgPosition.Add(boids[otherBoidId].position)
					avgVelocity = avgVelocity.Add(boids[otherBoidId].velocity)
					separation = separation.Add(b.position.Subtract(boids[otherBoidId].position).DivideByVal(distance))
				}
			}
		}
	}
	rWlock.RUnlock()

	accel := Vector2D{
		x: b.borderBounce(b.position.x, screenWidth),
		y: b.borderBounce(b.position.y, screenHeight),
	}

	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivideByVal(count), avgVelocity.DivideByVal(count)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyByVal(adjRate)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyByVal(adjRate)
		accelSeparation := separation.MultiplyByVal(adjRate)
		accel = accel.Add(accelCohesion).Add(accelAlignment).Add(accelSeparation)
	}

	return accel
}

func (b *Boid) moveOne() {
	// calcAcceleration should be outside lock that's below to avoid locking itself
	acceleration := b.calcAcceleration()

	rWlock.Lock()
	// limit method makes velocity no more than one pixel at a time so animation is smooth
	b.velocity = b.velocity.Add(acceleration).Limit(-1, 1)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	//next := b.position.Add(b.velocity)
	// boid bounces off of screen edge
	//if next.x >= screenWidth || next.x < 0 {
	//	b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	//}
	//if next.y >= screenHeight || next.y < 0 {
	//	b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	//}
	rWlock.Unlock()
}

func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(id int) {
	b := Boid{
		id:       id,
		position: Vector2D{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity: Vector2D{x: rand.Float64()*2 - 1, y: rand.Float64()*2 - 1}, // random number between -1 and 1
	}
	boids[id] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = id
	go b.start()
}
