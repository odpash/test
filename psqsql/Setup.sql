
CREATE TABLE items
(
    id integer,
    prices real[],
    sizes text[] COLLATE pg_catalog."default",
    colors text[] COLLATE pg_catalog."default",
    imagelinks text[] COLLATE pg_catalog."default",
    saleprices real[],
    infodate text[] COLLATE pg_catalog."default",
    count integer,
    category text COLLATE pg_catalog."default",
    items_id_key integer
);
