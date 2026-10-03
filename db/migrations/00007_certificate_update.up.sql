CREATE TABLE certificate_templates (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    language VARCHAR(5) NOT NULL,
    version VARCHAR(50) NOT NULL,
    layout_config JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE certificates (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    cert_hash VARCHAR(255) UNIQUE NOT NULL,
    file_key VARCHAR(255) UNIQUE NOT NULL,
    template_id BIGINT NOT NULL,
    issuer_id BIGINT NOT NULL,
    
    recipient_name VARCHAR(255) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    issued_date DATE NOT NULL,

    recipient_user_id INT NULL,
    event_id INT NULL,
    
    CONSTRAINT certificates_template_FK FOREIGN KEY (template_id) REFERENCES certificate_templates(id),
    CONSTRAINT certificates_user_FK FOREIGN KEY (recipient_user_id) REFERENCES users(id),
    CONSTRAINT certificates_event_FK FOREIGN KEY (event_id) REFERENCES event(id),
    INDEX idx_event (event_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;