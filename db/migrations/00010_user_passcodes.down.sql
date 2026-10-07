DROP TABLE users_passcodes;

CREATE TABLE users_passcodes (
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