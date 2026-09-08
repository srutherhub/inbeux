CREATE TYPE action_type AS ENUM (
    'create_reminder',
    'get_all_reminders',
    'delete_reminder',
    'update_user_timezone',
    'update_user_communication_preference',
    'update_user_name',
    'unknown'
);

CREATE TABLE actions (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  type action_type NOT NULL,
  message_id BIGINT NOT NULL REFERENCES messages(id),
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
)