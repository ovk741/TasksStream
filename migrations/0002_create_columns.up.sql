CREATE TABLE columns (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    board_id TEXT NOT NULL,
    position INT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_columns_board
        FOREIGN KEY (board_id)
        REFERENCES boards(id)
        ON DELETE CASCADE

    CONSTRAINT uniq_columns_board_position
        UNIQUE (board_id, position)    
);

