# Pulsar

Provide a brief description of what your application does and its unique features. You might also want to include badges for build status, test coverage, and other relevant indicators.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
  - [General](#general)
  - [Authentication](#authentication)
  - [Profile Management](#profile-management)
  - [Friendship Management](#friendship-management)
  - [Post Management](#post-management)
  - [Feed](#feed)

## Installation

To install the dependencies, run the following command:

```sh
go get ./...
```

## Configuration

Explain how to configure the application, e.g., environment variables, config files, etc.

## Usage

How to run the application:

```sh
go run main.go
```

## API Endpoints

Descriptions of the API endpoints and their functionalities.

### General

- `GET /api/ping` - Health check endpoint that returns `pong`.

### Authentication

- `POST /api/auth/sign-in` - Endpoint for user authentication.
- `POST /api/auth/register` - Endpoint for user registration.

### Profile Management

- `GET /api/me/profile` - Retrieve the authenticated user's profile.
- `PATCH /api/me/profile` - Update the authenticated user's profile.
- `POST /api/me/updatePassword` - Update the authenticated user's password.

### Friendship Management

- `POST /api/friends/add` - Add a friend for the authenticated user.
- `POST /api/friends/remove` - Remove a friend for the authenticated user.
- `GET /api/friends` - Get the friend list of the authenticated user.

### Post Management

- `POST /api/posts/new` - Create a new post.
- `POST /api/posts/:postId/like` - Like a post.
- `POST /api/posts/:postId/dislike` - Dislike a post.
- `GET /api/posts/:postId` - Get a post by ID.

### Feed

- `GET /api/posts/feed/:login` - Get the feed for a specific user.
- `GET /api/posts/feed/my` - Get the feed for the authenticated user.
