CREATE UNLOGGED TABLE authors
(
    id   bigserial not null,
    name text      not null default '',
    PRIMARY KEY (id)
);

CREATE INDEX authors_by_name ON authors (name);

CREATE UNLOGGED TABLE articles
(
    id           bigserial                not null,
    title        text                     not null default '',
    body         text                     not null default '',
    published_at timestamp with time zone not null,
    PRIMARY KEY (id)
);

CREATE INDEX articles_by_published_at ON articles (published_at);
