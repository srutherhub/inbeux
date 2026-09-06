CREATE TYPE messages_status AS ENUM ('pending', 'complete', 'classified', 'failed');
CREATE TYPE messages_source AS ENUM ('email');

CREATE TABLE messages (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  message_id VARCHAR(255) NOT NULL UNIQUE,
  source messages_source NOT NULL DEFAULT 'email',
  status messages_status NOT NULL DEFAULT 'pending',
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_status ON messages(status);