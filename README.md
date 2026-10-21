# Ranked Matchmaking Service

A small backend prototype built with Go and Redis that simulates
score-based matchmaking for online games.

## Features

- **Ranked queue:** Stores players and their scores in a Redis Sorted Set.
- **Concurrent queue joins:** Uses Go goroutines to simulate multiple
  players joining the matchmaking queue concurrently.
- **Room formation:** Groups queued players into rooms of five.
- **Queue management:** Removes matched players from the queue and
  stores room membership in Redis lists.

## Tech Stack

- Go
- Redis
- go-redis/v9
- Goroutines and sync.WaitGroup

## Prerequisites

- Go installed
- Redis running locally on port 6379

## Getting Started

1. Start your local Redis server.
2. Install the Go dependencies:

   ```bash
   go mod tidy