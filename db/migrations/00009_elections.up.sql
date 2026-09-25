CREATE TABLE config (
    `key` VARCHAR(255) PRIMARY KEY,
    `value` JSON NOT NULL
);

CREATE TABLE candidates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    cycle INT NOT NULL,
    belonging VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY unique_user_per_cycle (user_id, cycle),
    INDEX idx_candidate_cycle (cycle),

    CONSTRAINT fk_candidates_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE ticket_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    cycle INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_cycle_user UNIQUE (user_id, cycle),

    CONSTRAINT fk_ticket_records_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE vote_tickets (
    code VARCHAR(255) PRIMARY KEY,
    used BOOLEAN DEFAULT FALSE NOT NULL,
    used_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_ticket_used (used)
);

CREATE TABLE votes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    candidate_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_vote_candidate (candidate_id),
    INDEX idx_vote_created_at (created_at),

    CONSTRAINT fk_votes_candidate FOREIGN KEY (candidate_id) 
        REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE TABLE election_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    cycle INT NOT NULL,
    place INT NOT NULL,
    number_of_votes INT NOT NULL DEFAULT 0,
    belonging VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_results_cycle (cycle),
    INDEX idx_results_place (cycle, place),

    CONSTRAINT fk_election_results_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);