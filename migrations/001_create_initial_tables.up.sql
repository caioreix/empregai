CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TYPE ROLE AS ENUM ('admin', 'costumer', 'worker');

CREATE TYPE ADSTATUS AS ENUM ('opened', 'completed', 'canceled', 'expired');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(64) NOT NULL CHECK (email <> ''),
    password VARCHAR(250) NOT NULL CHECK (octet_length(password) <> 0),
    role ROLE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(email, role)
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(64) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(64) UNIQUE NOT NULL,
    description TEXT,
    type VARCHAR(64) NOT NULL,
    options JSON,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE advertisement_field (
    advertisement_id UUID,
    field_id UUID,
    value VARCHAR(250),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(advertisement_id, field_id)
);

CREATE TABLE advertisements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    costumer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    description TEXT NOT NULL,
    status ADSTATUS NOT NULL,
    currency CHAR(3) NOT NULL,
    price BIGINT,
    expiration_date TIMESTAMP WITH TIME ZONE NOT NULL,
    selected_cadidate UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE candidates (
    advertisement_id UUID,
    worker_id UUID,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (advertisement_id, worker_id)
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    advertisement_id UUID NOT NULL REFERENCES advertisements(id) ON DELETE CASCADE,
    rating DECIMAL,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
