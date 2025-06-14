CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE user_chats (
    user_id INT NOT NULL,
    chat_with_user_id INT NOT NULL,
    PRIMARY KEY (user_id, chat_with_user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_with_user_id) REFERENCES users(id)
);