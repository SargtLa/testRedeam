-- Title
-- Author
-- Publisher
-- Publish Date
-- Rating (1-3)
-- Status (CheckedIn, CheckedOut)

CREATE TABLE books
(
    id SERIAL NOT NULL,
    title VARCHAR NOT NULL,
    author VARCHAR NOT NULL,
    publisher VARCHAR(255) NOT NULL,
    publish_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    rating INTEGER NOT NULL DEFAULT 0,
    status BOOLEAN DEFAULT false,
    PRIMARY KEY (id)
);

COMMENT ON TABLE books IS 'list of books';

CREATE UNIQUE INDEX books_title_author_uindex
    on books (title, author);