import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";

export const IMPORT_ALBUM_URL = "/api/v1/music/album"

export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  importAlbum(formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData("/api/v1/music/album", "POST", api.PostAlbum, z.undefined(), formData, options)
  }
}
