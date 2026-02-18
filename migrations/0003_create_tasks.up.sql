CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    column_id TEXT NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_tasks_column
        FOREIGN KEY (column_id)
        REFERENCES columns(id)
        ON DELETE CASCADE

    CONSTRAINT uniq_task_column_position
        UNIQUE (column_id, position)       
);

