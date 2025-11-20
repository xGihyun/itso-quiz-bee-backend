#!/bin/sh

# Script to promote a user to admin role
# Usage: ./make-admin.sh <username>

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <username>"
  echo "Example: $0 john_doe"
  exit 1
fi

USERNAME="$1"

# Use the DATABASE_URL from environment, or construct it from individual vars
if [ -z "$DATABASE_URL" ]; then
  DB_USER="${DB_USER:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-postgres}"
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_NAME="${DB_NAME:-itso_quiz_bee}"
  DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
fi

echo "Attempting to promote user '$USERNAME' to admin..."
echo "Database URL: $DATABASE_URL"

# Execute the update query
psql "$DATABASE_URL" <<EOF
UPDATE users 
SET role = 'admin' 
WHERE username = '$USERNAME';

-- Show the result
SELECT user_id, username, role, created_at FROM users WHERE username = '$USERNAME';
EOF

echo "User promotion complete!"
