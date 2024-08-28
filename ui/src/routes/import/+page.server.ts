import { PostAlbumBody } from "$lib/api/types";
import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const formData = await request.formData();
    console.log(formData);

    const albumName = formData.get("albumName");
    if (!albumName) {
      throw error(500, "'albumName' not set");
    }

    const artistName = formData.get("artistName");
    if (!artistName) {
      throw error(500, "'artistName' not set");
    }

    const data: PostAlbumBody = {
      name: albumName.toString(),
      artist: artistName.toString(),
    };

    const body = new FormData();
    body.set("data", JSON.stringify(data));

    const files = formData.getAll("files");
    files.forEach((f) => {
      const file = f as File;
      body.append("files", file);
    });

    console.log(body);

    const res = await locals.apiClient.importAlbum(body);
    console.log(res);

    if (res.success) {
      throw redirect(301, `/edit/album/${res.data.id}`);
    }

    // TODO(patrik): Send error to page
  },
};
