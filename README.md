# Meeting Management API

A RESTful API built with **Go (Gin framework)** for managing users, meetings, and invitations.  
This project includes authentication, scheduling validation, and email notifications.

---

## ✨ Features

- **Users**

  - Register, login, profile management
  - JWT-based authentication

- **Meetings**

  - Create, edit, list, and delete meetings
  - Time conflict validation (no overlapping schedules)

- **Invitations**

  - Send and receive invitations
  - Email notifications for invitations
  - Accept or decline meeting invites

- **Dashboard**
  - List all meetings of the logged-in user

---

## 🗄️ Database Schema

**Users Table**

- `id` (PK)
- `name`
- `email`
- `password` (hashed)
- `created_at`

**Meetings Table**

- `id` (PK)
- `title`
- `description`
- `start_time`
- `end_time`
- `owner_id` (FK → users.id)
- `created_at`

**Invitations Table**

- `id` (PK)
- `meeting_id` (FK → meetings.id)
- `user_id` (FK → users.id)
- `status` (pending / accepted / declined)
- `created_at`

---

## 🛠️ Technologies

- **Backend:** Go + Gin
- **Database:** PostgreSQL
- **Authentication:** JWT
- **Email Notifications:** SMTP

---

## ⚙️ Setup

### Clone the repository

```bash
git clone https://github.com/jhonnydsl/gerenciamento-de-reunioes.git
cd gerenciamento-de-reunioes

# Create a .env file in the project root:

## Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASS=your_password
DB_NAME=meeting_db

## SMTP (for email invitations)
SMTP_HOST=smtp.yourprovider.com
SMTP_PORT=587
SMTP_USER=youremail@example.com
SMTP_PASS=yourpassword

## JWT
JWT_SECRET=your_secret_key
```
