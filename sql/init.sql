DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed) VALUES
                                                      ('Primary task', 'Learn a REST API', false),
                                                      ('Secondary task', 'Create a REST API program', false),
                                                      ('Task #3', 'Download Docker Desktop', true),
                                                      ('Task #4', 'Put it on GitHub', false),
                                                      ('Additional task', 'Clear some disk space up', true)