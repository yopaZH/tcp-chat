type User struct {
    Id uint64           // в базе — это primary key (уникальный идентификатор)
    Name string         // в базе — username
    Conn net.Conn       // не хранится в базе, это сетевое соединение, которое живёт только во время сессии
    ChatsWith map[uint64]struct{}  // множество id пользователей, с кем ведёт чат — нужно хранить отдельно
}

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