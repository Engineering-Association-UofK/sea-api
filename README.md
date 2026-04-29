# SEA Hub Backend

This is the Go-based backend for the steering engineering association's platform (at UofK). It handles everything from the public-facing CMS and event scheduling to internal tools like certificate generation and our bots.

## How to run it

We have two ways to get this running locally.

### 1. The Full Setup
Use this if you need to work on anything involving file upload or storage layer.

* **What it does:** Spins up the Go app along with a MySQL 8.4 and SeaweedFS services.
* **Setup:** 
1. `cp .env.compose.example .env.compose`
2. `docker-compose up --build`

### 2. The Lean Setup
Use this if you're just working on API logic or database workflows and don't want to run the whole container stack.

* **What this does:** It runs the app directly on your machine.
* **Setup:**
1. `cp .env.example .env` (Make sure you point it to a MySQL Compatible instance).
2. `go run cmd/api/main.go`

* **Note:** Any endpoint that works with the storage service will throw a **500** error in this mode because the file service won't be reachable.

## Project Features
The logic is inside `/internal/services`. what's in there:

* **Content:** CMS (e.g news, blogs, etc), media storage, and bot interactions.
* **Users:** There is Account management, user profiles, and stateless JWT auth.
* **Tools:** For PDF generation for certificates, event tracking, and collaborators management.
* **Infrastructure:** Task scheduler, SMTP mailer and a rate limiting.

## The Database
We use **MySQL 8.4** and **sqlx**.

* **Migrations:** There is no need for manual SQL imports. We use a migration library that runs on startup. Drop your `.sql` migration files in the `/db/migrations` folder.
* **No ORMs:** We do not use GORM or other wrappers, we write raw SQL queries. It keeps things fast and makes it clear what's happening.

