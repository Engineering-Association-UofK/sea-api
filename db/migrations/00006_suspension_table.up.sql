CREATE TABLE suspensions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL UNIQUE,
    reason TEXT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(idx)
);

CREATE TABLE suspension_history (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    reason TEXT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(idx)
);