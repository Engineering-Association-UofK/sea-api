CREATE TABLE events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    background_id INT NOT NULL,
    
    belonging VARCHAR(255) NOT NULL,

    require_applying TINYINT NOT NULL,
    max_applications INT NULL,
    
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

CREATE TABLE event_coords (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,

    FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE
);

CREATE TABLE event_participants (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    joined_at DATE NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE
);

CREATE TABLE event_application (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    Accepted TINYINT NOT NULL,
    started_at DATE NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE
);