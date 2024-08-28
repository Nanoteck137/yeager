import { env } from "$env/dynamic/private";
import { ApiClient } from "$lib/api/client";
import { redirect, type Handle } from "@sveltejs/kit";

const apiAddress = env.API_ADDRESS ? env.API_ADDRESS : "";

export const handle: Handle = async ({ event, resolve }) => {
  const url = new URL(event.request.url);

  let addr = apiAddress;
  if (addr == "") {
    addr = url.origin;
  }
  const client = new ApiClient(addr);
  event.locals.apiClient = client;

  if (url.pathname === "/login" && event.locals.user) {
    throw redirect(303, "/");
  }

  const response = await resolve(event);
  return response;
};
