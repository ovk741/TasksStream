CREATE TABLE board_members (
    id TEXT PRIMARY KEY,
    board_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_board_members_board
        FOREIGN KEY (board_id)
        REFERENCES boards(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_board_members_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uniq_board_user
        UNIQUE (board_id, user_id)
);