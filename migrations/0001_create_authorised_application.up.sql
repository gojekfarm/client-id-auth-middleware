CREATE TABLE authorized_applications (
  client_id VARCHAR(255) PRIMARY KEY ,
  pass_key VARCHAR(255) NOT NULL ,
  updated_at TIMESTAMP,
  created_at TIMESTAMP
)
