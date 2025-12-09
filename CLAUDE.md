# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- `npm run dev` - Start development server with Tailwind CSS compilation and nodemon auto-reload
- `npm run tailwind:css` - Compile Tailwind CSS from source to output file
- `npm start` - Start production server (node server.js)

## Architecture Overview

CannaNote is an Express.js web application for cannabis experience journaling with the following structure:

**Core Architecture:**
- **MVC Pattern**: Controllers handle business logic, models define data schemas, views render EJS templates
- **Session-based Authentication**: Uses express-session with MongoDB store (connect-mongo)
- **Database**: MongoDB with Mongoose ODM
- **Frontend**: EJS templating with Tailwind CSS for styling
- **Middleware**: Custom middleware in `middleware/mid.js` handles common setup

**Route Structure:**
- `/` - Unauthenticated routes (login, signup, marketing pages)
- `/cannanote` - Main application routes (entries CRUD) - requires authentication
- `/htmx` - HTMX-specific routes for dynamic interactions

**Key Models:**
- `Entry`: Cannabis experience entries with strain, type, amount, consumption method, description, and metadata (votes, favorites)
- `User`: User accounts for authentication

**Authentication Flow:**
- Unauthenticated users see marketing/login pages via `routes/unauth.js`
- All `/cannanote` routes require session authentication via middleware check
- Sessions stored in MongoDB with automatic cleanup

**Frontend Architecture:**
- EJS views in `views/` with reusable partials in `views/partials/`
- Tailwind CSS compilation from `public/css/tailwind.css` to `public/css/style.css`
- Static assets served from `public/` directory
- HTMX integration for dynamic interactions

## Environment Setup

Required environment variables:
- `DATABASE_URL` - MongoDB connection string
- `SECRET` - Session secret key
- `PORT` - Server port (defaults to 3001)

## Database Schema

**Entry Schema** (`models/entries.js`):
- username, strain, type, amount, consumption, description
- date (auto-generated), tags array
- meta object with votes and favorites count

## Styling System

- Tailwind CSS configured to scan `views/**/*.ejs` files
- Source file: `public/css/tailwind.css`
- Output file: `public/css/style.css` (compiled via npm script)
- PostCSS configuration in `postcss.config.js`