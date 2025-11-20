#!/bin/sh

# Script to import quiz questions from a CSV file and auto-create the quiz
# Usage: ./import-questions.sh <csv_file>
# CSV Format: quiz_name,question_number,question_text,question_type,points,duration_seconds,answer1,answer2,answer3,answer4

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <csv_file>"
  echo "Example: $0 questions.csv"
  echo ""
  echo "CSV Format:"
  echo "quiz_name,question_number,question_text,question_type,points,duration_seconds,answer1,answer2,answer3,answer4"
  echo ""
  echo "- question_type: 'multiple-choice', 'boolean', or 'written'"
  echo "- duration_seconds: Time limit per question in seconds (e.g., 60)"
  echo "- For 'written' questions, leave answer fields empty"
  echo "- The quiz will be auto-created if it doesn't exist"
  exit 1
fi

CSV_FILE="$1"

if [ ! -f "$CSV_FILE" ]; then
  echo "Error: File '$CSV_FILE' not found!"
  exit 1
fi

# Use the DATABASE_URL from environment, or construct it from individual vars
if [ -z "$DATABASE_URL" ]; then
  DB_USER="${DB_USER:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-postgres}"
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_NAME="${DB_NAME:-itso_quiz_bee}"
  DATABASE_URL="postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
fi

echo "Importing questions from '$CSV_FILE'..."
echo "Database: $DATABASE_URL"
echo ""

# Counter for tracking
QUESTION_COUNT=0
ANSWER_COUNT=0
CURRENT_QUIZ_NAME=""
CURRENT_QUIZ_ID=""

# Read CSV file and process each line (skip header)
tail -n +2 "$CSV_FILE" | while IFS=',' read -r QUIZ_NAME QUESTION_NUM QUESTION_TEXT QUESTION_TYPE POINTS DURATION ANSWER1 ANSWER2 ANSWER3 ANSWER4; do
  
  # Trim whitespace
  QUIZ_NAME=$(echo "$QUIZ_NAME" | xargs)
  QUESTION_NUM=$(echo "$QUESTION_NUM" | xargs)
  QUESTION_TEXT=$(echo "$QUESTION_TEXT" | xargs)
  QUESTION_TYPE=$(echo "$QUESTION_TYPE" | xargs)
  POINTS=$(echo "$POINTS" | xargs)
  DURATION=$(echo "$DURATION" | xargs)
  ANSWER1=$(echo "$ANSWER1" | xargs)
  ANSWER2=$(echo "$ANSWER2" | xargs)
  ANSWER3=$(echo "$ANSWER3" | xargs)
  ANSWER4=$(echo "$ANSWER4" | xargs)
  CORRECT_ANSWER=$(echo "$CORRECT_ANSWER" | xargs)

  # Skip empty lines
  if [ -z "$QUIZ_NAME" ]; then
    continue
  fi

  # Create or retrieve the quiz if the name changed
  if [ "$QUIZ_NAME" != "$CURRENT_QUIZ_NAME" ]; then
    echo ""
    echo "Processing quiz: $QUIZ_NAME"
    
    # Check if quiz already exists
    EXISTING_QUIZ=$(psql "$DATABASE_URL" -t -A -c "SELECT quiz_id FROM quizzes WHERE name = E'$QUIZ_NAME' LIMIT 1;")
    
    if [ -n "$EXISTING_QUIZ" ]; then
      CURRENT_QUIZ_ID="$EXISTING_QUIZ"
      echo "  → Found existing quiz (ID: $CURRENT_QUIZ_ID)"
    else
      # Create new quiz (no lobby_id needed anymore)
      CURRENT_QUIZ_ID=$(psql "$DATABASE_URL" -t -A -c "
        INSERT INTO quizzes (name, status)
        VALUES (E'$QUIZ_NAME', 'open')
        RETURNING quiz_id;
      ")
      echo "  → Created new quiz (ID: $CURRENT_QUIZ_ID)"
    fi
    
    CURRENT_QUIZ_NAME="$QUIZ_NAME"
  fi

  echo "  Q$QUESTION_NUM - $QUESTION_TEXT (Type: $QUESTION_TYPE)"

  # Insert question and capture the ID
  QUESTION_ID=$(psql "$DATABASE_URL" -t -A -c "
    INSERT INTO quiz_questions (quiz_id, content, points, order_number, duration)
    VALUES ('$CURRENT_QUIZ_ID', E'$QUESTION_TEXT', $POINTS, $QUESTION_NUM, '${DURATION} seconds'::interval)
    RETURNING quiz_question_id;
  ")

  echo "  → Question ID: $QUESTION_ID"

  # Handle different question types
  if [ "$QUESTION_TYPE" = "multiple-choice" ]; then
    # Insert multiple choice answers (no is_correct column anymore)
    psql "$DATABASE_URL" > /dev/null <<EOF
      INSERT INTO quiz_answers (quiz_question_id, content)
      VALUES 
        ('$QUESTION_ID', E'$ANSWER1'),
        ('$QUESTION_ID', E'$ANSWER2'),
        ('$QUESTION_ID', E'$ANSWER3'),
        ('$QUESTION_ID', E'$ANSWER4');
EOF
    echo "    → Added 4 answers"
    ANSWER_COUNT=$((ANSWER_COUNT + 4))

  elif [ "$QUESTION_TYPE" = "boolean" ]; then
    # Insert boolean answers
    psql "$DATABASE_URL" > /dev/null <<EOF
      INSERT INTO quiz_answers (quiz_question_id, content)
      VALUES 
        ('$QUESTION_ID', 'true'),
        ('$QUESTION_ID', 'false');
EOF
    echo "    → Added 2 answers"
    ANSWER_COUNT=$((ANSWER_COUNT + 2))

  elif [ "$QUESTION_TYPE" = "written" ]; then
    echo "    → Written answer (no preset answers)"
  fi

  QUESTION_COUNT=$((QUESTION_COUNT + 1))
done

echo ""
echo "✓ Import complete!"
echo "  - Questions imported: $QUESTION_COUNT"
echo "  - Answers imported: $ANSWER_COUNT"
