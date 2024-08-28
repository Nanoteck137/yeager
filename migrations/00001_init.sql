-- +goose Up
CREATE TABLE artists (
    id TEXT,
    name TEXT NOT NULL,
    picture TEXT,

    CONSTRAINT artists_pk PRIMARY KEY (id)
);

CREATE TABLE albums (
    id TEXT,
    artist_id TEXT NOT NULL,

    name TEXT NOT NULL,
    cover_art TEXT,

    CONSTRAINT albums_pk PRIMARY KEY(id),

    CONSTRAINT albums_artist_id_fk FOREIGN KEY (artist_id)
        REFERENCES artists(id)
);

CREATE TABLE tracks (
    id TEXT NOT NULL,

    album_id TEXT NOT NULL,
    artist_id TEXT NOT NULL,

    title TEXT NOT NULL,

    cover_art TEXT,

    track_number INT,
    duration INT,
    year INT,

    original_file TEXT NOT NULL,
    transcoded_file TEXT NOT NULL,

    CONSTRAINT tracks_pk PRIMARY KEY(album_id, id),

    CONSTRAINT tracks_album_id_fk FOREIGN KEY (album_id)
        REFERENCES albums(id),

    CONSTRAINT tracks_artist_id_fk FOREIGN KEY (artist_id)
        REFERENCES artists(id)
);


-- +goose Down
DROP TABLE tracks;
DROP TABLE albums;
DROP TABLE artists;
