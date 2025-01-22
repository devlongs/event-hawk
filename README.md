# EventHawk 🦅

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A real-time Ethereum smart contract event monitor with customizable alerts, built in Go. Never miss important on-chain activity!


## Features ✨

- 📡 **Real-time Event Monitoring** via WebSocket
- 🔍 Filter events by contract address and signatures
- 🔔 **Slack Notifications** for detected events
- 📝 **Text/JSON Output** formats
- ⚙️ YAML Configuration
- 🔄 Basic Chain Reorganization (Reorg) Handling

## Installation ⚡

### Prerequisites
- Go 1.20+ installed
- Ethereum node access (e.g., [Infura](https://infura.io/) WebSocket URL)
- Slack webhook URL (optional)

### Build from Source
```bash
git clone https://github.com/devlongs/eventhawk.git
cd eventhawk
go mod download
go build -o eventhawk .