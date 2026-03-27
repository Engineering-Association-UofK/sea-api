ALTER TABLE user_roles ADD CONSTRAINT unique_user_role UNIQUE (user_id, role);

ALTER TABLE users ADD COLUMN profile_image_id BIGINT;
ALTER TABLE users ADD CONSTRAINT fk_profile_image FOREIGN KEY (profile_image_id) REFERENCES files(id);

-- New files system

CREATE TABLE files (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    file_key TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(255) NOT NULL
);

-- Gallery as a subtype file to make it easier to manage

CREATE TABLE gallery_assets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    file_id BIGINT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    alt_text VARCHAR(255) NOT NULL,
    uploaded_by BIGINT NOT NULL,
    showcase TINYINT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_uploaded_by (uploaded_by),
    INDEX idx_file_id (file_id),
    
    FOREIGN KEY (file_id) REFERENCES files(id)
);

CREATE TABLE gallery_references (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    object_type TINYINT NOT NULL,
    object_id BIGINT NOT NULL,

    UNIQUE KEY unique_asset_object (asset_id, object_type, object_id),

    INDEX idx_object (object_type, object_id),
    INDEX idx_asset (asset_id),
    
    FOREIGN KEY (asset_id) REFERENCES gallery_assets(id)
);

-- CMS

CREATE TABLE blog_posts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content LONGTEXT NOT NULL,
    cover_image_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    is_published BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_published_date (is_published, created_at),
    INDEX idx_cover_image (cover_image_id),

    FOREIGN KEY (cover_image_id) REFERENCES gallery_assets(id)
);

CREATE TABLE team_members (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    role VARCHAR(100) NOT NULL,
    bio TEXT NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order (is_active, display_order)
);

CREATE TABLE news_gallery (
    id INT AUTO_INCREMENT PRIMARY KEY,
    asset_id BIGINT NOT NULL,

    INDEX idx_asset_id (asset_id),
    
    FOREIGN KEY (asset_id) REFERENCES gallery_assets(id)
);

-- Forms engine schema

CREATE TABLE forms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    allow_multiple BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    header_image_id BIGINT, 
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_is_active (is_active),

    FOREIGN KEY (header_image_id) REFERENCES gallery_assets(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE form_pages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    form_id BIGINT NOT NULL,
    page_num INT NOT NULL,

    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE,
    UNIQUE KEY unique_form_page (form_id, page_num)
);

CREATE TABLE form_questions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    form_page_id BIGINT NOT NULL,
    question_text TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    options JSON NOT NULL,
    is_required BOOLEAN DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    FOREIGN KEY (form_page_id) REFERENCES form_pages(id) ON DELETE CASCADE,
    INDEX idx_page_order (form_page_id, display_order)
);

CREATE TABLE form_responses (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    form_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status ENUM('DRAFT', 'SUBMITTED') DEFAULT 'DRAFT',
    submitted_at DATETIME NULL,
    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE,
    INDEX idx_user_form (user_id, form_id)
);

CREATE TABLE form_answers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    response_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,
    answer_value TEXT NOT NULL,
    FOREIGN KEY (response_id) REFERENCES form_responses(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES form_questions(id) ON DELETE CASCADE,

    UNIQUE KEY unique_response_question (response_id, question_id),
    
    INDEX idx_question_answer (question_id, answer_value(100)) 
);