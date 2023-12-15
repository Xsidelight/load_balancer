# Simple Go Load Balancer

This project implements a basic load balancer in Go. It's designed to distribute incoming TCP connections among a set of backend servers in a round-robin fashion. Additionally, the load balancer periodically checks the health of each backend server to ensure reliability and availability.

## Features

- **Round-Robin Load Balancing**: Distributes incoming connections among backend servers in a round-robin manner.
- **Health Check**: Periodically checks the health of each backend server and logs their status.
- **Concurrency and Synchronization**: Uses Go's concurrency features and the `sync` package for efficient and safe operation.

## How It Works

The load balancer listens for incoming TCP connections on a specified port. Each incoming connection is forwarded to one of the backend servers based on a round-robin algorithm. Concurrently, the load balancer performs health checks on each backend server every 10 seconds, logging their status.

The implementation uses Go's `net`, `http`, and `sync` packages. The `sync.Once` type is particularly used to ensure that the health check goroutines for the backend servers are started only once, even when called from multiple goroutines.

Disclaimer
This project is developed as a part of a coding challenge from [Coding Challenges FYI](https://codingchallenges.fyi/challenges/challenge-load-balancer). It is intended for educational purposes and not recommended for production use without further modifications and testing.
