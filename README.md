# BHIV Blockchain Inspector - Verification Tool

A comprehensive blockchain analysis, debugging, and monitoring tool built for the BHIV blockchain network. This tool provides error detection, multi-node consensus analysis, and real-time monitoring capabilities.

## 🎯 Project Overview

This inspector was developed over 3 days with progressive feature additions:

- **Day 1**: Core functionality (loading, scanning, comparing, RPC support, watch mode)
- **Day 2**: Multi-node consensus analysis and fork detection
- **Day 3**: Production packaging, reporting, and comprehensive testing

## ✨ Features

### Day 1 - Core Functionality

#### 1. **Data Loading**
- Load sample blockchain data into LevelDB
- Generate test blocks with valid hashes
- Configurable block count

#### 2. **Block Viewing**
- View individual blocks by height
- Support for both LevelDB and RPC sources
- JSON and human-readable output formats

#### 3. **Comprehensive Error Scanning**
Detects 11 types of blockchain errors:
- Corrupted JSON blocks
- Invalid hash computations
- Timestamp anomalies (future, past, non-increasing)
- Duplicate hashes
- Empty blocks
- Broken chain linkage (prevHash errors)
- Height inconsistencies
- Missing blocks
- Out-of-order blocks

#### 4. **Two-Node Comparison**
- Compare blockchain state between two nodes
- Identify divergence points
- Calculate sync percentage
- Generate repair recommendations

#### 5. **RPC Support**
- Fetch blocks from live nodes via HTTP
- Health endpoint monitoring
- Configurable RPC endpoints

#### 6. **Real-Time Watch Mode**
- Continuous node monitoring
- Customizable polling intervals
- Visual status indicators
- Latency tracking

### Day 2 - Consensus Analysis

#### 1. **Multi-Node Analysis**
- Compare 3+ nodes simultaneously
- Config-based node management
- Comprehensive network state analysis

#### 2. **Fork Detection**
- Identify blockchain forks
- Track affected nodes per fork
- Analyze fork branches

#### 3. **Canonical Chain Identification**
- Determine the "correct" chain
- Score-based chain selection
- Agreement-based consensus

#### 4. **Network Health Assessment**
- Calculate network health score
- Grade network status (Excellent/Good/Fair/Poor/Critical)
- Identify synchronization issues

#### 5. **Auto-Repair Recommendations**
- Generate actionable recommendations
- Identify nodes needing resync
- Prioritize critical issues

### Day 3 - Production Features

#### 1. **Comprehensive Reporting**
- JSON report generation
- Combined analysis data
- Timestamped snapshots

#### 2. **Multi-Platform Builds**
- Windows (amd64)
- Linux (amd64)
- macOS (amd64, arm64)

#### 3. **Build Automation**
- Makefile support
- Batch scripts for Windows
- One-command builds

## 📦 Installation

### Prerequisites

- Go 1.19 or higher
- Git (optional, for cloning)

### Quick Install

