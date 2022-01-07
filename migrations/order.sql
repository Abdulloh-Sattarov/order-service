create  database orders;
drop table orders;

create  table orders (
    order_id UUID PRIMARY KEY,
    book_uuid UUID,
    description VARCHAR(256),
    created_at timestamp DEFAULT current_timestamp,
    updated_at timestamp  DEFAULT  current_timestamp,
    deleted_at timestamp DEFAULT null
);

insert into orders (order_id, book_uuid, author_id, description) values ('e740486e-2810-40fc-a122-7ef01a1a7684', 'e740486e-2810-40fc-a122-7ef01a1a7684', 'e740486e-2810-40fc-a122-7ef01a1a7684', 'ok');