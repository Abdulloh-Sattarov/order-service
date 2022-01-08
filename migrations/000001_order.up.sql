create  table orders (
     order_id UUID PRIMARY KEY,
     book_uuid UUID,
     description VARCHAR(256),
     created_at timestamp DEFAULT current_timestamp,
     updated_at timestamp  DEFAULT  current_timestamp,
     deleted_at timestamp DEFAULT null
);