import { z } from "zod";

export const FormSchema = z.object({
  albumName: z.string().min(1),
  artistName: z.string().min(1),
  year: z.preprocess((a) => {
    if (a === "") return undefined;
    return parseInt(a as string, 10);
  }, z.number().positive().optional()),
});

export type Form = z.infer<typeof FormSchema>;
export type FormErrors = z.inferFlattenedErrors<
  typeof FormSchema
>["fieldErrors"];

export const NewAlbumServerSchema = z.object({
  albumName: z.string().min(1),
  artistName: z.string().min(1),
  year: z.number().optional(),
  tracks: z.array(
    z.object({
      title: z.string().min(1),
      number: z.number().positive().optional(),
    }),
  ),
});

export type NewAlbumServer = z.infer<typeof NewAlbumServerSchema>;
