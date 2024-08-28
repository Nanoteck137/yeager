package types

import "path"

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d WorkDir) ArtistsDir() string {
	return path.Join(d.String(), "artists")
}

func (d WorkDir) AlbumsDir() string {
	return path.Join(d.String(), "albums")
}

func (d WorkDir) GeneratedDir() string {
	return path.Join(d.String(), "generated")
}

func (d WorkDir) GeneratedArtistsDir() string {
	return path.Join(d.GeneratedDir(), "artists")
}

func (d WorkDir) GeneratedAlbumsDir() string {
	return path.Join(d.GeneratedDir(), "albums")
}
