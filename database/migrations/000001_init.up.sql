-- Enable UUID extension for unique identifiers
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Category Table
CREATE TABLE category (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256) UNIQUE NOT NULL,
    parent_id UUID NULL,  -- Supports nested categories
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (parent_id) REFERENCES category(id) ON DELETE CASCADE
);

CREATE INDEX idx_category_parent ON category(parent_id);

-- Links Table
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT UNIQUE NOT NULL,  -- Ensures no duplicate links
    title VARCHAR(256) NOT NULL,
    description TEXT NOT NULL,
    short_url TEXT UNIQUE NOT NULL,  -- For shortened URLs
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_links_url ON links(url);
CREATE INDEX idx_links_shorturl ON links(short_url);

-- Link-Category Association Table
CREATE TABLE link_category_map (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id UUID NOT NULL,  -- References the link
    category_id UUID NOT NULL,  -- References the category
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE,
    UNIQUE(link_id, category_id)
);

CREATE INDEX idx_link_category_map_link ON link_category_map(link_id);
CREATE INDEX idx_link_category_map_category ON link_category_map(category_id); 