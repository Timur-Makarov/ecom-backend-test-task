CREATE TABLE banners
(
    id   INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE counter_statistics
(
    banner_id      INTEGER NOT NULL,
    timestamp_from BIGINT  NOT NULL,
    timestamp_to   BIGINT  NOT NULL,
    count          BIGINT  NOT NULL,
    PRIMARY KEY (timestamp_from, timestamp_to)
);