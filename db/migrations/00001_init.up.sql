-- -------- FILE SYSTEM SCHEMA

CREATE TABLE files (
    id INT AUTO_INCREMENT PRIMARY KEY,
    file_key TEXT NOT NULL,
    file_size INT NOT NULL,
    mime_type VARCHAR(255) NOT NULL
);

-- -------- USER SCHEMA

CREATE TABLE users (
    id INT PRIMARY KEY,
    uni_id INT UNIQUE,
    username VARCHAR(255) UNIQUE NOT NULL,

    name_ar VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,

    gender VARCHAR(10) NOT NULL,
    department VARCHAR(255) NOT NULL,
    profile_image_id INT,

    password VARCHAR(255) NOT NULL,
    verified tinyint NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    UNIQUE (profile_image_id),
    FOREIGN KEY (profile_image_id) REFERENCES files(id)
);

CREATE TABLE users_temp (
    id INT PRIMARY KEY,
    uni_id INT UNIQUE,
    username VARCHAR(255) UNIQUE,

    name_ar VARCHAR(255),
    name_en VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,

    password VARCHAR(255),
    verified tinyint NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE user_roles (
    user_id INT NOT NULL,
    role VARCHAR(50) NOT NULL,
    INDEX (user_id),

    UNIQUE (user_id, role),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE
);

CREATE TABLE verification_code (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    code TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE
);

-- -------- EVENT SCHEMA

CREATE TABLE event (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    presenter_id INT NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    max_participants INT NOT NULL,
    outcomes TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    FOREIGN KEY (presenter_id) REFERENCES collaborators(id) ON UPDATE CASCADE
);

CREATE TABLE event_component (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    max_score DECIMAL(10, 2) NOT NULL,

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE event_participant (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    grade DECIMAL(5, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    joined_at DATE NOT NULL,
    completed TINYINT NOT NULL,

    UNIQUE (event_id, user_id),

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE component_score (
    id INT PRIMARY KEY AUTO_INCREMENT,
    participant_id INT NOT NULL,
    component_id INT NOT NULL,
    score DECIMAL(10, 2),

    UNIQUE (participant_id, component_id),

    FOREIGN KEY (participant_id) REFERENCES event_participant(id) ON DELETE CASCADE,
    FOREIGN KEY (component_id) REFERENCES event_component(id) ON DELETE CASCADE
);

CREATE TABLE collaborators (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name_ar VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    signature_id INT NOT NULL,

    FOREIGN KEY (signature_id) REFERENCES files(id) ON DELETE CASCADE
)

-- -------- CERTIFICATE SCHEMA

CREATE TABLE certificate (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cert_hash VARCHAR(255) UNIQUE NOT NULL,
    user_id INT NOT NULL,
    event_id INT NOT NULL,
    grade DECIMAL(5, 2) NOT NULL,
    issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL,

    UNIQUE (user_id, event_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (event_id) REFERENCES event(id)
);

CREATE TABLE certificate_file (
    id INT PRIMARY KEY AUTO_INCREMENT,
    certificate_id INT NOT NULL,
    store_id INT NOT NULL UNIQUE,
    lang VARCHAR(10) NOT NULL,

    UNIQUE (certificate_id, lang),

    FOREIGN KEY (certificate_id) REFERENCES certificate(id),
    FOREIGN KEY (store_id) REFERENCES files(id) ON DELETE CASCADE
);

-- -------- SUSPENSION SCHEMA

CREATE TABLE suspensions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL UNIQUE,
    reason TEXT NOT NULL,
    admin_id INT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NOT NULL,

    FOREIGN KEY (admin_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE
);

CREATE TABLE suspension_history (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    reason TEXT NOT NULL,
    admin_id INT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NOT NULL,

    FOREIGN KEY (admin_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE
);

-- -------- GALLERY SCHEMA

CREATE TABLE gallery_assets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    file_id INT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    alt_text VARCHAR(255) NOT NULL,
    uploaded_by INT NOT NULL,
    showcase TINYINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_uploaded_by (uploaded_by),
    INDEX idx_file_id (file_id),
    
    FOREIGN KEY (uploaded_by) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id)
);

CREATE TABLE gallery_references (
    id INT AUTO_INCREMENT PRIMARY KEY,
    asset_id INT NOT NULL,
    object_type TINYINT NOT NULL,
    object_id INT NOT NULL,

    UNIQUE KEY unique_asset_object (asset_id, object_type, object_id),

    INDEX idx_object (object_type, object_id),
    INDEX idx_asset (asset_id),
    
    FOREIGN KEY (asset_id) REFERENCES gallery_assets(id)
);

-- -------- CMS SCHEMA

CREATE TABLE blog_posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content LONGTEXT NOT NULL,
    cover_image_id INT NOT NULL,
    author_id INT NOT NULL,
    is_published BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_published_date (is_published, created_at),
    INDEX idx_cover_image (cover_image_id),

    FOREIGN KEY (author_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (cover_image_id) REFERENCES gallery_assets(id)
);

CREATE TABLE team_members (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    role VARCHAR(100) NOT NULL,
    bio TEXT NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order (is_active, display_order),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE news_gallery (
    id INT AUTO_INCREMENT PRIMARY KEY,
    asset_id INT NOT NULL,

    INDEX idx_asset_id (asset_id),
    
    FOREIGN KEY (asset_id) REFERENCES gallery_assets(id)
);

-- -------- FORMS SCHEMA

CREATE TABLE forms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    allow_multiple BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    header_image_id INT, 
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_is_active (is_active),

    FOREIGN KEY (header_image_id) REFERENCES gallery_assets(id),
    FOREIGN KEY (created_by) REFERENCES users(id) ON UPDATE CASCADE
);

CREATE TABLE form_pages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    form_id INT NOT NULL,
    page_num INT NOT NULL,

    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE,
    UNIQUE KEY unique_form_page (form_id, page_num)
);

CREATE TABLE form_questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    form_page_id INT NOT NULL,
    question_text TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    options JSON NOT NULL,
    is_required BOOLEAN DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    FOREIGN KEY (form_page_id) REFERENCES form_pages(id) ON DELETE CASCADE,
    INDEX idx_page_order (form_page_id, display_order)
);

CREATE TABLE form_responses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    form_id INT NOT NULL,
    user_id INT NOT NULL,
    status ENUM('DRAFT', 'SUBMITTED') DEFAULT 'DRAFT',
    submitted_at DATETIME NULL,

    INDEX idx_user_form (user_id, form_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
);

CREATE TABLE form_answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    response_id INT NOT NULL,
    question_id INT NOT NULL,
    answer_value TEXT NOT NULL,
    
    INDEX idx_question_answer (question_id, answer_value(100)),

    FOREIGN KEY (response_id) REFERENCES form_responses(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES form_questions(id) ON DELETE CASCADE
);

-- RATE LIMIT SCHEMA

CREATE TABLE rate_limits (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    request_count INT DEFAULT 1,
    last_request TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE (ip_address, endpoint)
);