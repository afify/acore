# Acore Project — Developer Setup & Deployment

Welcome to **acore**! This document walks you through setting up a local development environment and configuring GitHub for zero-downtime EC2 deployments, without diving into code snippets.

---

## 📋 Prerequisites

Make sure you have the following tools installed on your machine:

- Git
- Go (version 1.21 or later)
- Docker & Docker Compose
- Make (GNU Make)
- AWS CLI (optional, for EC2 key management)

---

## 🏠 Local Environment Setup

1. **Clone the repository**
   Use your preferred Git client or the command line to get the `acore` code onto your local machine.

2. **Create and populate your `.env` file**
   Duplicate the example environment file in the repo root, then fill in values for things like application name, OAuth credentials, database URL, and any other service keys you need for local testing.

3. **Start supporting services**
   Use Docker Compose to launch the database, cache, and any other infrastructure containers required by the app.

4. **Build and run the application**
   Rely on the provided Makefile targets (e.g. build, migrate, run) to compile the Go code and start the server. Once running, you’ll be able to access the app on the configured localhost port.

---

## 🔑 Environment Variables

Your local `.env` should include values such as:

- Application settings (name, secrets, callback URLs)
- Database connection string
- Third-party API keys

Ensure you never commit real secrets—use placeholders in the example file.

---

## 🛡️ GitHub Secrets for EC2 Deployment

To enable automated, blue/green deploys to EC2, add the following **four** repository secrets under **Settings → Secrets and variables → Actions**:

1. **EC2_USER**
   The SSH username on your EC2 instance (for example, `ec2-user` or `ubuntu`).

2. **EC2_HOST**
   The public DNS name or IP address of your EC2 server.

3. **EC2_SSH_PRIVATE_KEY**
   The private key (PEM format) that lets GitHub Actions SSH into the EC2 host.

4. **EC2_GIT_SSH_PRIVATE_KEY**
   A deploy-only SSH key that allows the EC2 instance to pull from your Git repo securely.

---

## 🚀 Deployment Workflow Overview

- **Push to `main`**
  Triggers the GitHub Actions pipeline.

- **Checkout using the deploy key**
  Actions clone the repo on the runner with minimal permissions.

- **SSH into EC2**
  Uses your EC2 SSH key to connect to the server.

- **Blue/green swap**
  The pipeline pulls the latest code on the inactive “color”, builds the new container, then switches traffic over with zero downtime.

- **Post-deploy health check**
  Confirms the new version is serving correctly before retiring the old one.

---

## 🎯 Summary

1. Install all prerequisites locally
2. Clone the repo and configure your `.env`
3. Bring up Docker infrastructure and start the app via Make
4. In GitHub, add the four EC2-related secrets.
5. Push to `main` and watch the automated EC2 deployment execute.

You’re now ready to work on **acore** and safely deploy updates to production with confidence!
